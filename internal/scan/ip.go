package scan

import (
	"net"
	"time"
)

func ScanRange(ips []string, port string) []string {
	var active []string

	return active
}

func Ping(address string) bool {
	conn, err := net.DialTimeout("ip4:icmp", address, time.Second)
	if err == nil {
		defer conn.Close()
		return true
	}

	time.Sleep(time.Second) // Wait for goroutines to complete
	return false
}
