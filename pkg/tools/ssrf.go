package tools

import (
	"fmt"
	"net"
	"net/url"
)

// checkSSRF 校验出站请求地址, 阻止指向服务自身、云元数据、链路本地等危险地址的 SSRF。
//
// 说明: 为保留自建内部监控(内网 Prometheus/Loki 等)的正当用途, 不拦截 RFC1918 私网段,
// 仅拦截回环(127.0.0.0/8、::1)、链路本地(169.254.0.0/16, 含云元数据 169.254.169.254)、
// 组播与未指定地址; 域名解析失败时放行, 交由 http client 自身的解析兜底。
func checkSSRF(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("非法 URL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("出站请求仅允许 http/https 协议, 当前: %s", u.Scheme)
	}

	host := u.Hostname()
	if ip := net.ParseIP(host); ip != nil {
		if isBlockedIP(ip) {
			return fmt.Errorf("禁止访问的地址: %s", host)
		}
		return nil
	}

	// 域名: 解析后逐项校验
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil
	}
	for _, ip := range ips {
		if isBlockedIP(ip) {
			return fmt.Errorf("禁止访问的地址: %s (%s)", host, ip)
		}
	}
	return nil
}

// isBlockedIP 判断 IP 是否属于禁止出站访问的地址范围。
func isBlockedIP(ip net.IP) bool {
	return ip.IsLoopback() ||
		ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() ||
		ip.IsMulticast() ||
		ip.IsUnspecified()
}
