package dnslb

import (
	"fmt"
	"log"
	"os"

	"github.com/cloudflare/cloudflare-go"
)

// Cleanup removes unhealthy endpoints.
func Cleanup(subdomain string, ports []string) error {
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

	records, err := api.DNSRecords(zoneID, cloudflare.DNSRecord{
		Name: subdomain + "." + zone,
	})
	if err != nil {
		return fmt.Errorf("load dns records: %v", err)
	}

	for _, record := range records {
		log.Printf("record %v %v", record.Type, record.Content)

		switch record.Type {
		case "A":
			ipv4 := &ipv4{ip: ip{addr: record.Content}}

			err := ipv4.cleanup(api, zoneID, domain, ports)
			if err != nil {
				log.Printf("cleanup record %s: %v", record.Content, err)
			}
		case "AAAA":
			ipv6 := &ipv6{ip: ip{addr: record.Content}}

			err := ipv6.cleanup(api, zoneID, domain, ports)
			if err != nil {
				log.Printf("cleanup record %s: %v", record.Content, err)
			}
		}
	}

	return nil
}
