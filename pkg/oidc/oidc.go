package oidc

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"io"
	"net/http"
	"net/url"
	"strings"
	"watchAlert/internal/types"
)

func GetOauthToken(upper string, code string) (*types.OauthToken, error) {
	ctx := context.Background()

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)

	host := fmt.Sprintf("%s/v1/oauth/token", upper)

	request, err := http.NewRequest(http.MethodPost, host, strings.NewReader(form.Encode()))
	if err != nil {
		logc.Error(ctx, fmt.Sprintf("post 请求建立失败, err: %s", err.Error()))
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(request)
	if err != nil {
		logc.Error(ctx, fmt.Sprintf("post 请求发送失败, err: %s", err.Error()))
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logc.Error(ctx, fmt.Sprintf("读取响应体失败, err: %s", err.Error()))
		return nil, err
	}

	var result types.OauthToken
	if err = json.Unmarshal(body, &result); err != nil {
		logc.Error(ctx, fmt.Sprintf("解析响应体失败, err: %s", err.Error()))
		return nil, err
	}

	if result.AccessToken == "" || result.RefreshToken == "" {
		logc.Error(ctx, fmt.Sprintf("请求 accessToken 或 refreshToken 失败, body: %s", string(body)))
		return nil, fmt.Errorf("failed to get oauth token")
	}

	return &result, nil
}

func GetCurrentUser(upper string, uid interface{}) (*types.RespOidcUserInfo, error) {
	ctx := context.Background()

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	host := fmt.Sprintf("%s/v1/users/%v?app_label=oidc", upper, uid)

	request, err := http.NewRequest(http.MethodGet, host, nil)
	if err != nil {
		logc.Error(ctx, fmt.Sprintf("post 请求建立失败, err: %s", err.Error()))
		return nil, err
	}

	resp, err := client.Do(request)
	if err != nil {
		logc.Error(ctx, fmt.Sprintf("post 请求发送失败, err: %s", err.Error()))
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logc.Error(ctx, fmt.Sprintf("读取响应体失败, err: %s", err.Error()))
		return nil, err
	}

	var result types.RespOidcUserInfo
	if err = json.Unmarshal(body, &result); err != nil {
		logc.Error(ctx, fmt.Sprintf("解析响应体失败, err: %s", err.Error()))
		return nil, err
	}

	if result.Code != 0 {
		logc.Error(ctx, fmt.Sprintf("获取用户信息失败, 错误码: %d, body: %s", result.Code, string(body)))
		return nil, fmt.Errorf("failed to get user info, code: %d", result.Code)
	}

	return &result, nil
}

func DecodeToken(upper, accessToken, rfToken string) error {
	ctx := context.Background()

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	host := fmt.Sprintf("%s/v1/auth", upper)

	request, err := http.NewRequest(http.MethodGet, host, nil)
	if err != nil {
		logc.Error(ctx, fmt.Sprintf("post 请求建立失败, err: %s", err.Error()))
		return err
	}

	request.Header.Set("Authorization", accessToken)
	request.Header.Set("RFToken", rfToken)

	resp, err := client.Do(request)
	if err != nil {
		logc.Error(ctx, fmt.Sprintf("post 请求发送失败, err: %s", err.Error()))
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logc.Error(ctx, fmt.Sprintf("读取响应体失败, err: %s", err.Error()))
		return err
	}

	var result types.RespOidcAuth
	if err = json.Unmarshal(body, &result); err != nil {
		logc.Error(ctx, fmt.Sprintf("解析响应体失败, err: %s", err.Error()))
		return err
	}

	if result.Code != 0 {
		logc.Error(ctx, fmt.Sprintf("认证失败, 错误码: %d", result.Code))
		return fmt.Errorf("auth failed: %d", result.Code)
	}

	return nil
}
