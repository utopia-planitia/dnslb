package dnslb

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Endpoint runs the life cycle of a endpoint.
func Endpoint(ports []string, ipv4Enabled, ipv6Enabled bool) error {
	api, err := initAPI(os.Getenv("CF_API_TOKEN"), os.Getenv("CF_API_KEY"), os.Getenv("CF_API_EMAIL"))
	if err != nil {
		return fmt.Errorf("init api: %v", err)
	}

	log.Println("connected api")

	zone, zoneID, err := loadZoneID(api)
	if err != nil {
		return fmt.Errorf("load zone: %v", err)
	}

	log.Printf("connected zone: %v", zone)

	subdomain := os.Getenv("CF_SUBDOMAIN")
	if subdomain == "" {
		return fmt.Errorf("CF_SUBDOMAIN needs to be set")
	}

	log.Printf("subdomain: %v", subdomain)

	domain := subdomain + "." + zone
	log.Printf("domain: %v", domain)

	log.Printf("ports: %v", ports)

	ipv4, ipv6, err := detectIPs(ipv4Enabled, ipv6Enabled)
	if err != nil {
		return fmt.Errorf("detect local IPs: %v", err)
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		if ipv4Enabled {
			err := ipv4.maintain(api, zoneID, domain, ports)
			if err != nil {
				log.Println(err)
			}
		}

		if ipv6Enabled {
			err = ipv6.maintain(api, zoneID, domain, ports)
			if err != nil {
				log.Println(err)
			}
		}

		select {
		case <-sigs:
			log.Println("shutting down")

			if ipv4Enabled {
				err = ipv4.remove(api, zoneID, domain)
				if err != nil {
					log.Printf("remove IPv4 %v: %v", ipv4.addr, err)
				}
			}

			if ipv6Enabled {
				err = ipv6.remove(api, zoneID, domain)
				if err != nil {
					log.Printf("remove IPv6 %v: %v", ipv4.addr, err)
				}
			}

			log.Println("await DNS TTL (120 seconds)")
			time.Sleep(120 * time.Second)

			return nil
		case <-ticker.C:
			continue
		}
	}
}

func detectIPs(ipv4Enabled, ipv6Enabled bool) (*ipv4, *ipv6, error) {
	if !ipv4Enabled && !ipv6Enabled {
		return nil, nil, fmt.Errorf("at least one of --ipv4 and --ipv6 need to be enabled")
	}

	var addr4 *ipv4

	var addr6 *ipv6

	if ipv4Enabled {
		addr, err := myIP("tcp4")
		if err != nil {
			return nil, nil, fmt.Errorf("lookup IPv4 address: %v", err)
		}

		log.Printf("ipv4: %v", addr)

		addr4 = &ipv4{ip: ip{addr: addr}}
	}

	if ipv6Enabled {
		addr, err := myIP("tcp6")
		if err != nil {
			return nil, nil, fmt.Errorf("lookup IPv6 address: %v", err)
		}

		log.Printf("ipv6: %v", addr)

		addr6 = &ipv6{ip: ip{addr: addr}}
	}

	return addr4, addr6, nil
}
