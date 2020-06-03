package dnslb

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

// IPv4 print the detected local IPv4 address to stdout.
func IPv4() error {
	ip, err := myIP("tcp4")
	if err != nil {
		return err
	}

	_, err = io.WriteString(os.Stdout, ip+"\n")
	if err != nil {
		return err
	}

	return nil
}

// IPv6 print the detected local IPv6 address to stdout.
func IPv6() error {
	ip, err := myIP("tcp6")
	if err != nil {
		return err
	}

	_, err = io.WriteString(os.Stdout, ip)
	if err != nil {
		return err
	}

	return nil
}

// Based on https://community.cloudflare.com/t/can-1-1-1-1-be-used-to-find-out-ones-public-ip-address/14971/5.
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
