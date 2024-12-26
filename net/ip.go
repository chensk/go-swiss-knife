package net

import (
	"errors"
	"net"
)

func GetCurrentIpv4() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("unable to get ip")
}

func GetCurrentIpv6() string {
	// google dns
	conn, err := net.Dial("udp", "[2001:4860:4860::8888]:80")
	if err != nil {
		return ""
	}
	defer func() {
		_ = conn.Close()
	}()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
