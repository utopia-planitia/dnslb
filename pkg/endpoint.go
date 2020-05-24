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
func Endpoint(subdomain string, ports []string, IPv4, IPv6 bool) error {
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

	domain := subdomain + "." + zone
	log.Printf("domain: %v", domain)

	log.Printf("ports: %v", ports)

	if !IPv4 && !IPv6 {
		return fmt.Errorf("at least one of --ipv4 and --ipv6 need to be enabled")
	}

	ipv4, err := myIPv4(IPv4)
	if err != nil {
		return fmt.Errorf("lookup IPv4 address: %v", err)
	}

	ipv6, err := myIPv6(IPv6)
	if err != nil {
		return fmt.Errorf("lookup IPv6 address: %v", err)
	}

	if ipv4 != nil {
		log.Printf("ipv4: %v", ipv4.addr)
	}

	if ipv6 != nil {
		log.Printf("ipv6: %v", ipv4.addr)
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		if ipv4 != nil {
			err := ipv4.maintain(api, zoneID, domain, ports)
			if err != nil {
				log.Println(err)
			}
		}

		if ipv6 != nil {
			err = ipv6.maintain(api, zoneID, domain, ports)
			if err != nil {
				log.Println(err)
			}
		}

		select {
		case <-sigs:
			log.Println("shutting down")

			if ipv4 != nil {
				err = ipv4.remove(api, zoneID, domain)
				if err != nil {
					log.Printf("remove IPv4 %v: %v", ipv4.addr, err)
				}
			}

			if ipv6 != nil {
				err = ipv6.remove(api, zoneID, domain)
				if err != nil {
					log.Printf("remove IPv6 %v: %v", ipv4.addr, err)
				}
			}

			log.Println("await DNS TTL (120 seconds)")
			time.Sleep(120 * time.Second)

			return nil
		case <-ticker.C:
		}
	}
}

func myIPv4(enabled bool) (*ipv4, error) {
	if !enabled {
		return nil, nil
	}

	ip4, err := myIP("tcp4")
	if err != nil {
		return nil, err
	}

	return &ipv4{ip: ip{addr: ip4}}, nil
}

func myIPv6(enabled bool) (*ipv6, error) {
	if !enabled {
		return nil, nil
	}

	ip6, err := myIP("tcp6")
	if err != nil {
		return nil, err
	}

	return &ipv6{ip: ip{addr: ip6}}, nil
}
