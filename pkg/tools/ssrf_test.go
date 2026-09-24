package tools

import (
	"net"
	"testing"
)

func TestCheckSSRF(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		wantErr bool
	}{
		{"公网 https", "https://example.com/path", false},
		{"IPv4 回环", "http://127.0.0.1:9090", true},
		{"IPv6 回环", "http://[::1]:9090", true},
		{"云元数据地址", "http://169.254.169.254/latest/meta-data", true},
		{"链路本地", "http://169.254.10.10/api", true},
		{"内网私网(放行, 供内部监控)", "http://10.0.0.1:9090", false},
		{"私网 B 段(放行)", "http://172.16.1.2/api", false},
		{"私网 C 段(放行)", "http://192.168.1.1/api", false},
		{"非 http 协议", "ftp://example.com/file", true},
		{"非法 URL", "://bad url", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkSSRF(tt.rawURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkSSRF(%q) err = %v, wantErr = %v", tt.rawURL, err, tt.wantErr)
			}
		})
	}
}

func TestIsBlockedIP(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		blocked bool
	}{
		{"127.0.0.1", "127.0.0.1", true},
		{"127.8.8.8", "127.8.8.8", true},
		{"::1", "::1", true},
		{"169.254.169.254", "169.254.169.254", true},
		{"224.0.0.1 组播", "224.0.0.1", true},
		{"0.0.0.0", "0.0.0.0", true},
		{"公网 8.8.8.8", "8.8.8.8", false},
		{"私网 10.0.0.1", "10.0.0.1", false},
		{"私网 172.20.0.1", "172.20.0.1", false},
		{"私网 192.168.0.1", "192.168.0.1", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isBlockedIP(net.ParseIP(tt.ip)); got != tt.blocked {
				t.Errorf("isBlockedIP(%q) = %v, want %v", tt.ip, got, tt.blocked)
			}
		})
	}
}
