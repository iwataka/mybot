package main

import (
	"net"
)

// getGlobalIP returns the IP network via which this computer is accessible.
func getGlobalIP() (*net.IP, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	maskSize := 0
	var result *net.IP
	for _, a := range as {
		ip, ipnet, err := net.ParseCIDR(a.String())
		if err != nil {
			return nil, err
		}
		if !ip.IsLoopback() && ip.IsGlobalUnicast() {
			size, _ := ipnet.Mask.Size()
			if size > maskSize {
				maskSize = size
				result = &ip
			}
		}
	}
	return result, nil
}
