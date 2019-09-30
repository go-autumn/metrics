// Author SunJun <i@sjis.me>
package util

import (
	"net"
	"strconv"
	"strings"
)

func IntranetIP() []string {
	ips := make([]string, 0)
	ices, _ := net.Interfaces()

	for _, ice := range ices {
		if ice.Flags&net.FlagUp == 0 {
			continue // interface down
		}

		if ice.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		// ignore docker and warden bridge
		if strings.HasPrefix(ice.Name, "docker") || strings.HasPrefix(ice.Name, "w-") {
			continue
		}

		adds, _ := ice.Addrs()
		for _, addr := range adds {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}

			ipStr := ip.String()
			if IsIntranet(ipStr) {
				ips = append(ips, ipStr)
			}
		}
	}

	return ips
}

func IsIntranet(ipStr string) bool {
	if strings.HasPrefix(ipStr, "10.") || strings.HasPrefix(ipStr, "192.168.") {
		return true
	}

	if strings.HasPrefix(ipStr, "172.") {
		// 172.16.0.0-172.31.255.255
		arr := strings.Split(ipStr, ".")
		if len(arr) != 4 {
			return false
		}

		second, err := strconv.ParseInt(arr[1], 10, 64)
		if err != nil {
			return false
		}

		if second >= 16 && second <= 31 {
			return true
		}
	}

	return false
}
