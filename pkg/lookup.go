package dnslb

import (
	"context"
	"fmt"
	"net"
	"time"
)

func myIPs(IPv4, IPv6 bool) ([]string, error) {
	if !IPv4 && !IPv6 {
		return []string{}, fmt.Errorf("at least one of IPv4 and IPv6 need to be selected")
	}

	ips := []string{}

	if IPv4 {
		ipv4, err := myIP("tcp4")
		if err != nil {
			return []string{}, fmt.Errorf("lookup IPv4 address: %v", err)
		}

		ips = append(ips, ipv4)
	}

	if IPv6 {
		ipv6, err := myIP("tcp6")
		if err != nil {
			return []string{}, fmt.Errorf("lookup IPv6 address: %v", err)
		}

		ips = append(ips, ipv6)
	}

	return ips, nil
}

// based on https://community.cloudflare.com/t/can-1-1-1-1-be-used-to-find-out-ones-public-ip-address/14971/5
// dig -6 TXT +short o-o.myaddr.l.google.com @ns1.google.com
// dig -4 TXT +short o-o.myaddr.l.google.com @ns1.google.com
func myIP(network string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r := &net.Resolver{
		Dial: func(ctx context.Context, _, _ string) (net.Conn, error) {
			h := net.JoinHostPort("ns1.google.com", "53")
			d := net.Dialer{}
			return d.DialContext(ctx, network, h)
		},
	}

	addrs, err := r.LookupTXT(ctx, "o-o.myaddr.l.google.com")
	if err != nil {
		return "", err
	}

	if len(addrs) == 0 {
		return "", fmt.Errorf("no IP address found")
	}

	if len(addrs) > 1 {
		return "", fmt.Errorf("too many IP addresses found")
	}

	return addrs[0], nil
}
