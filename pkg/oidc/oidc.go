package oidc

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"watchAlert/internal/types"
	"watchAlert/pkg/tools"
)

// wellKnownOidcConfiguration 是 OIDC Discovery 规范定义的标准路径。
const wellKnownOidcConfiguration = "/.well-known/openid-configuration"

func GetOpenIDConfiguration(upper string) (*types.RespOpenIDConfiguration, error) {
	resp, err := tools.Get(nil, resolveDiscoveryURL(upper), 10)
	if err != nil {
		return nil, err
	}

	var d types.RespOpenIDConfiguration
	if err = tools.ParseReaderBody(resp.Body, &d); err != nil {
		return nil, err
	}

	return &d, nil
}

// resolveDiscoveryURL 将 Issuer 基础地址解析为 OIDC Discovery 标准路径。
// 前端 oidc-client 以 Issuer 为 authority 并自动拼接 /.well-known/openid-configuration,
// 这里保持一致; 同时兼容已把完整 discovery URL 存入 UpperURI 的历史配置。
func resolveDiscoveryURL(issuer string) string {
	if strings.HasSuffix(issuer, wellKnownOidcConfiguration) {
		return issuer
	}
	return strings.TrimSuffix(issuer, "/") + wellKnownOidcConfiguration
}

func GetOauthToken(tokenUrl, code, clientID, clientSecret string) (*types.OauthToken, error) {
	header := make(map[string]string)
	header["Content-Type"] = "application/x-www-form-urlencoded"

	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	if clientID != "" {
		form.Add("client_id", clientID)
	}
	if clientSecret != "" {
		form.Add("client_secret", clientSecret)
	}

	resp, err := tools.Post(header, tokenUrl, bytes.NewReader([]byte(form.Encode())), 10)
	if err != nil {
		return nil, err
	}

	var d types.OauthToken
	if err = tools.ParseReaderBody(resp.Body, &d); err != nil {
		return nil, err
	}

	if d.AccessToken == "" {
		return nil, fmt.Errorf("failed to get oauth token: empty access_token")
	}

	return &d, nil
}

func GetCurrentUser(userInfoUrl, token string) (*types.RespOidcUserInfo, error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + token

	resp, err := tools.Get(header, userInfoUrl, 10)
	if err != nil {
		return nil, err
	}

	var d types.RespOidcUserInfo
	if err = tools.ParseReaderBody(resp.Body, &d); err != nil {
		return nil, err
	}

	return &d, nil
}
