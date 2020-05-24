package dnslb

import (
	"fmt"
	"log"

	"github.com/cloudflare/cloudflare-go"
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

func (i ipv4) cleanup(api *cloudflare.API, zoneID, domain string, ports []string) error {
	if !i.healthy(ports) {
		err := i.remove(api, zoneID, domain)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i ipv6) cleanup(api *cloudflare.API, zoneID, domain string, ports []string) error {
	if !i.healthy(ports) {
		err := i.remove(api, zoneID, domain)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i ipv4) maintain(api *cloudflare.API, zoneID, domain string, ports []string) error {
	action := i.remove
	if i.healthy(ports) {
		action = i.add
	}

	err := action(api, zoneID, domain)
	if err != nil {
		return err
	}

	return nil
}

func (i ipv6) maintain(api *cloudflare.API, zoneID, domain string, ports []string) error {
	action := i.remove
	if i.healthy(ports) {
		action = i.add
	}

	err := action(api, zoneID, domain)
	if err != nil {
		return err
	}

	return nil
}

func (i ip) healthy(ports []string) bool {
	ok, _ := allHealthyTCP(i.addr, ports)
	return ok
}

func (i ipv4) add(api *cloudflare.API, zoneID, domain string) error {
	return ensureRecord(api, zoneID, domain, i.addr, "A")
}

func (i ipv6) add(api *cloudflare.API, zoneID, domain string) error {
	return ensureRecord(api, zoneID, domain, i.addr, "AAAA")
}

func (i ip) remove(api *cloudflare.API, zoneID, domain string) error {
	return removeRecord(api, zoneID, domain, i.addr)
}

func ensureRecord(api *cloudflare.API, zoneID, domain, ip, typee string) error {
	log.Printf("ensure IP: %v", ip)
	records, err := api.DNSRecords(zoneID, cloudflare.DNSRecord{
		Name:    domain,
		Content: ip,
	})
	if err != nil {
		return fmt.Errorf("lookup record %v/%v: %v", domain, ip, err)
	}

	if len(records) > 0 {
		return nil
	}

	_, err = api.CreateDNSRecord(zoneID, cloudflare.DNSRecord{
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

func removeRecord(api *cloudflare.API, zoneID, domain, ip string) error {
	log.Printf("remove IP: %v", ip)

	records, err := api.DNSRecords(zoneID, cloudflare.DNSRecord{
		Name:    domain,
		Content: ip,
	})
	if err != nil {
		return fmt.Errorf("identify record %v/%v: %v", domain, ip, err)
	}

	for _, record := range records {
		err = api.DeleteDNSRecord(zoneID, record.ID)
		if err != nil {
			return fmt.Errorf("delete record %v: %v", record.ID, err)
		}

		log.Println("entry deleted")
	}

	return nil
}
