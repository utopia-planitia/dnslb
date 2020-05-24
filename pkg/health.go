package dnslb

import (
	"net"
	"time"
)

func allHealthyTCP(host string, ports []string) (bool, error) {
	for _, port := range ports {
		ok, err := isHealthyTCP(host, port)
		if !ok {
			return false, err
		}
	}
	return true, nil
}

func isHealthyTCP(host, port string) (bool, error) {
	timeout := time.Second
	//	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", port), timeout)
	if err != nil {
		return false, err
	}

	if conn != nil {
		defer conn.Close()
	}

	return true, nil
}
