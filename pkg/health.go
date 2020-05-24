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
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), time.Second)
	if err != nil {
		return false, err
	}

	if conn != nil {
		defer conn.Close()
	}

	return true, nil
}
