package oidc

import (
	"bytes"
	"fmt"
	"net/url"
	"watchAlert/internal/types"
	"watchAlert/pkg/tools"
)

func GetOpenIDConfiguration(upper string) (*types.RespOpenIDConfiguration, error) {
	resp, err := tools.Get(nil, fmt.Sprintf("%s/.well-known/openid-configuration", upper), 10)
	if err != nil {
		return nil, err
	}

	var d types.RespOpenIDConfiguration
	if err = tools.ParseReaderBody(resp.Body, &d); err != nil {
		return nil, err
	}

	return &d, nil
}

func GetOauthToken(tokenUrl, code string) (*types.OauthToken, error) {
	header := make(map[string]string)
	header["Content-Type"] = "application/x-www-form-urlencoded"

	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)

	resp, err := tools.Post(header, tokenUrl, bytes.NewReader([]byte(form.Encode())), 10)
	if err != nil {
		return nil, err
	}

	var d types.OauthToken
	if err = tools.ParseReaderBody(resp.Body, &d); err != nil {
		return nil, err
	}

	if d.AccessToken == "" || d.RefreshToken == "" {
		return nil, fmt.Errorf("failed to get oauth token")
	}

	return &d, nil
}

func GetCurrentUser(userInfoUrl, token string) (*types.RespOidcUserInfo, error) {
	header := make(map[string]string)
	header["Authorization"] = token

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
