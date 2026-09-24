package oidc

import "testing"

func TestResolveDiscoveryURL(t *testing.T) {
	tests := []struct {
		name   string
		issuer string
		want   string
	}{
		{"issuer 基础地址", "https://idp.example.com", "https://idp.example.com/.well-known/openid-configuration"},
		{"issuer 带末尾斜杠", "https://idp.example.com/", "https://idp.example.com/.well-known/openid-configuration"},
		{"issuer 带路径前缀", "https://idp.example.com/casdoor", "https://idp.example.com/casdoor/.well-known/openid-configuration"},
		{"已是完整 discovery URL", "https://idp.example.com/.well-known/openid-configuration", "https://idp.example.com/.well-known/openid-configuration"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveDiscoveryURL(tt.issuer); got != tt.want {
				t.Errorf("resolveDiscoveryURL(%q) = %q, want %q", tt.issuer, got, tt.want)
			}
		})
	}
}
