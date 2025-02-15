package dnslb

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudflare/cloudflare-go/v4"
)

type ipv4 struct {
	ip
}

type ipv6 struct {
	ip
}

type ip struct {
	addr string
}

func (i ipv4) cleanup(ctx context.Context, api *cloudflare.API, zoneID, domain string, ports []string) error {
	if !i.healthy(ports) {
		err := i.remove(ctx, api, zoneID, domain)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i ipv6) cleanup(ctx context.Context, api *cloudflare.API, zoneID, domain string, ports []string) error {
	if !i.healthy(ports) {
		err := i.remove(ctx, api, zoneID, domain)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i ipv4) maintain(ctx context.Context, api *cloudflare.API, zoneID, domain string, ports []string) error {
	action := i.remove
	if i.healthy(ports) {
		action = i.add
	}

	err := action(ctx, api, zoneID, domain)
	if err != nil {
		return err
	}

	return nil
}

func (i ipv6) maintain(ctx context.Context, api *cloudflare.API, zoneID, domain string, ports []string) error {
	action := i.remove
	if i.healthy(ports) {
		action = i.add
	}

	err := action(ctx, api, zoneID, domain)
	if err != nil {
		return err
	}

	return nil
}

func (i ip) healthy(ports []string) bool {
	ok, _ := allHealthyTCP(i.addr, ports)
	return ok
}

func (i ipv4) add(ctx context.Context, api *cloudflare.API, zoneID, domain string) error {
	return ensureRecord(ctx, api, zoneID, domain, i.addr, "A")
}

func (i ipv6) add(ctx context.Context, api *cloudflare.API, zoneID, domain string) error {
	return ensureRecord(ctx, api, zoneID, domain, i.addr, "AAAA")
}

func (i ip) remove(ctx context.Context, api *cloudflare.API, zoneID, domain string) error {
	return removeRecord(ctx, api, zoneID, domain, i.addr)
}

func ensureRecord(ctx context.Context, api *cloudflare.API, zoneID, domain, ip, typee string) error {
	log.Printf("ensure IP: %v", ip)

	records, _, err := api.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(zoneID), cloudflare.ListDNSRecordsParams{
		Name:    domain,
		Content: ip,
	})
	if err != nil {
		return fmt.Errorf("lookup records %v/%v: %v", domain, ip, err)
	}

	if len(records) > 0 {
		return nil
	}

	_, err = api.CreateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID), cloudflare.CreateDNSRecordParams{
		Name:    domain,
		Content: ip,
		Type:    typee,
		TTL:     120,
	})
	if err != nil {
		return fmt.Errorf("create record %v/%v: %v", domain, ip, err)
	}

	log.Println("entry created")

	return nil
}

func removeRecord(ctx context.Context, api *cloudflare.API, zoneID, domain, ip string) error {
	log.Printf("remove IP: %v", ip)

	records, _, err := api.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(zoneID), cloudflare.ListDNSRecordsParams{
		Name:    domain,
		Content: ip,
	})
	if err != nil {
		return fmt.Errorf("identify record %v/%v: %v", domain, ip, err)
	}

	for _, record := range records {
		err = api.DeleteDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID), record.ID)
		if err != nil {
			return fmt.Errorf("delete record %v: %v", record.ID, err)
		}

		log.Println("entry deleted")
	}

	return nil
}
