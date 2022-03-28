package dohttp

import "net"

// IsPublicIP 判断指定 IPv4 是否为公网 IP
//
// 字符串类型的 IP 地址 可以通过 net.ParseIP() 函数转为 net.IP
//
// @see https://blog.csdn.net/whatday/article/details/109689258
func IsPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}
