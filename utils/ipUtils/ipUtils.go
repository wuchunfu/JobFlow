package main

import (
	"errors"
	"math"
	"net"
	"net/http"
	"strings"
)

// 获取本机网卡IP
func getLocalIP() (ipv4 string, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet // IP地址
		isIpNet bool
	)
	// 获取所有网卡
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}
	// 取第一个非loop的网卡IP
	for _, addr = range addrs {
		// 这个网络地址是IP地址: ipv4, ipv6
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			// 跳过IPV6
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String() // 192.168.1.1
				return
			}
		}
	}
	return
}

// https://blog.thinkeridea.com/201903/go/get_client_ip.html
// HasLocalIPAddr 检测 IP 地址字符串是否是内网地址
func HasLocalIPAddr(ip string) bool {
	return HasLocalIP(net.ParseIP(ip))
}

// HasLocalIP 检测 IP 地址是否是内网地址
// 通过直接对比ip段范围效率更高，详见：https://github.com/thinkeridea/go-extend/issues/2
func HasLocalIP(ip net.IP) bool {
	if ip.IsLoopback() {
		return true
	}
	ip4 := ip.To4()
	if ip4 == nil {
		return false
	}
	return ip4[0] == 10 || // 10.0.0.0/8
		(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) || // 172.16.0.0/12
		(ip4[0] == 169 && ip4[1] == 254) || // 169.254.0.0/16
		(ip4[0] == 192 && ip4[1] == 168) // 192.168.0.0/16
}

// ClientIP 尽最大努力实现获取客户端 IP 的算法。
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作。
func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

// ClientPublicIP 尽最大努力实现获取客户端公网 IP 的算法。
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作。
func ClientPublicIP(r *http.Request) string {
	var ip string
	for _, ip = range strings.Split(r.Header.Get("X-Forwarded-For"), ",") {
		ip = strings.TrimSpace(ip)
		if ip != "" && !HasLocalIPAddr(ip) {
			return ip
		}
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" && !HasLocalIPAddr(ip) {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		if !HasLocalIPAddr(ip) {
			return ip
		}
	}
	return ""
}

// RemoteIP 通过 RemoteAddr 获取 IP 地址， 只是一个快速解析方法。
func RemoteIP(r *http.Request) string {
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

// IPString2Long 把ip字符串转为数值
func IPString2Long(ip string) (uint, error) {
	b := net.ParseIP(ip).To4()
	if b == nil {
		return 0, errors.New("invalid ipv4 format")
	}
	return uint(b[3]) | uint(b[2])<<8 | uint(b[1])<<16 | uint(b[0])<<24, nil
}

// Long2IPString 把数值转为ip字符串
func Long2IPString(i uint) (string, error) {
	if i > math.MaxUint32 {
		return "", errors.New("beyond the scope of ipv4")
	}
	ip := make(net.IP, net.IPv4len)
	ip[0] = byte(i >> 24)
	ip[1] = byte(i >> 16)
	ip[2] = byte(i >> 8)
	ip[3] = byte(i)
	return ip.String(), nil
}

// IP2Long 把net.IP转为数值
func IP2Long(ip net.IP) (uint, error) {
	b := ip.To4()
	if b == nil {
		return 0, errors.New("invalid ipv4 format")
	}
	return uint(b[3]) | uint(b[2])<<8 | uint(b[1])<<16 | uint(b[0])<<24, nil
}

// Long2IP 把数值转为net.IP
func Long2IP(i uint) (net.IP, error) {
	if i > math.MaxUint32 {
		return nil, errors.New("beyond the scope of ipv4")
	}
	ip := make(net.IP, net.IPv4len)
	ip[0] = byte(i >> 24)
	ip[1] = byte(i >> 16)
	ip[2] = byte(i >> 8)
	ip[3] = byte(i)
	return ip, nil
}
