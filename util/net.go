package util

import (
	"net"
	"strings"
)

func GetLoopbackInterface() (*net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for i := range ifaces {
		iface := ifaces[i]

		// 过滤无效地址
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		// 过滤非环回地址
		if iface.Flags&net.FlagLoopback == 0 {
			continue
		}
		// 过滤虚拟地址
		name := iface.Name
		name = strings.ToLower(name)
		if strings.Contains(name, "vmnet") {
			continue
		}
		return &iface, nil
	}
	return nil, nil
}

func GetInterfaceIPv4(iface *net.Interface) string {
	addrs, err := iface.Addrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil {
				if ip := ipnet.IP.To4().String(); ip != "" {
					return ip
				}
			}
		}
	}
	return ""
}
