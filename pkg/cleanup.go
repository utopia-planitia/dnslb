package dnslb

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloudflare/cloudflare-go/v2"
)

// CleanupLoop executes a cleanup every delay
func CleanupLoop(ctx context.Context, ports []string, delay time.Duration) error {
	for {
		err := func() error {
			ctx, cancel := context.WithTimeout(ctx, time.Minute)
			defer cancel()

			err := Cleanup(ctx, ports)
			if err != nil {
				return fmt.Errorf("cleanup failed: %v", err)
			}

			return nil
		}()
		if err != nil {
			return err
		}

		time.Sleep(delay)
	}
}

// Cleanup removes unhealthy endpoints.
func Cleanup(ctx context.Context, ports []string) error {
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

	records, _, err := api.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(zoneID), cloudflare.ListDNSRecordsParams{
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

			err := ipv4.cleanup(ctx, api, zoneID, domain, ports)
			if err != nil {
				log.Printf("cleanup record %s: %v", record.Content, err)
			}
		case "AAAA":
			ipv6 := &ipv6{ip: ip{addr: record.Content}}

			err := ipv6.cleanup(ctx, api, zoneID, domain, ports)
			if err != nil {
				log.Printf("cleanup record %s: %v", record.Content, err)
			}
		}
	}

	return nil
}
