package utils

import (
	"net"
	"os"
	"strconv"
	"strings"

	"math/rand"
)

// 私有地址和保留地址的列表
var (
	privateAndReservedIPs = []string{
		"10.0.0.0/8",     // 私有地址
		"172.16.0.0/12",  // 私有地址
		"192.168.0.0/16", // 私有地址
		"127.0.0.0/8",    // 环回地址
		"169.254.0.0/16", // 链路本地地址
		"::1/128",        // IPv6环回地址
		"fc00::/7",       // IPv6私有地址
		"fe80::/10",      // IPv6链路本地地址
		"fec0::/10",
	}
)

func generateRandomIPv4() net.IP {
	return net.IP{
		byte(rand.Intn(256)),
		byte(rand.Intn(256)),
		byte(rand.Intn(256)),
		byte(rand.Intn(256)),
	}
}

// generateRandomIPv6 生成一个随机的IPv6地址
func generateRandomIPv6() net.IP {
	var ip [16]byte
	for i := range ip {
		ip[i] = byte(rand.Intn(256))
	}
	return ip[:]
}

// GeneratePublicIPs 生成一个包含指定数量的IPv4和IPv6公网IP地址的列表 ipType 1-IPv4  2-IPv6
func GeneratePublicIPs(count int, ipType uint8) ([]net.IP, error) {
	var publicIPs []net.IP
	for len(publicIPs) < count {
		if ipType == 1 {
			ip := generateRandomIPv4()
			if !isPrivateOrReserved(ip) {
				publicIPs = append(publicIPs, ip)
			}
		} else {
			ip := generateRandomIPv6()
			if !isPrivateOrReserved(ip) {
				publicIPs = append(publicIPs, ip)
			}
		}
	}
	return publicIPs, nil
}

// isPrivateOrReserved 检查IP是否是私有地址或保留地址
func isPrivateOrReserved(ip net.IP) bool {
	for _, cidr := range privateAndReservedIPs {
		_, ipnet, _ := net.ParseCIDR(cidr)
		if ipnet.Contains(ip) {
			return true
		}
	}
	return false
}

// ReadProxy 读取指定目录下的指定txt文件，并将每一行内容存入数组 proxyType 1-http 2-https 3-socks5
func ReadProxy(filePath string, proxyType uint8) ([]string, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(fileContent), "\n")
	if len(lines) > 0 {
		return removeErrLines(lines, proxyType), nil
	}
	return nil, nil
}

func removeErrLines(lines []string, proxyType uint8) []string {
	var result []string
	var proxyscheme string
	if proxyType == 1 {
		proxyscheme = "http://"
	} else if proxyType == 2 {
		proxyscheme = "https://"
	} else if proxyType == 3 {
		proxyscheme = "socks5://"
	}
	for _, line := range lines {
		if line != "" && checkIPPort(line) {
			result = append(result, proxyscheme+line)
		}
	}
	return result
}

func checkIPPort(s string) bool {
	host, port, err := net.SplitHostPort(s)
	if err != nil {
		return false
	}
	if host == "" || port == "" {
		return false
	}
	if net.ParseIP(host) == nil {
		return false
	}
	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		return false
	}
	return true
}
