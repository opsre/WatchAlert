package tools

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logc"
)

// defaultTransport 是共享的 HTTP Transport, 通过连接池复用 TCP/TLS 连接,
// 避免每次 Get/Post 都重新拨号与握手(此前每请求 new 一个 Transport, 连接池失效)。
// 注意: InsecureSkipVerify 为历史行为(统一跳过证书校验), 属待处理的独立安全项。
var defaultTransport = &http.Transport{
	TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	},
	Proxy:               http.ProxyFromEnvironment,
	MaxIdleConns:        100,
	MaxIdleConnsPerHost: 10,
	IdleConnTimeout:     90 * time.Second,
	DisableKeepAlives:   false,
}

// do 执行请求: 复用 defaultTransport 连接池, 用每请求 http.Client 的 Timeout 控制超时。
// 注意: 不能用 context.WithTimeout + defer cancel 包在返回 response 的外层 ——
// cancel 会在 do 返回时立即触发, 导致调用方后续 io.ReadAll(resp.Body) 读到
// "context canceled"(响应体较大、分片传输时必现)。
func do(request *http.Request, timeout int) (*http.Response, error) {
	client := &http.Client{
		Transport: defaultTransport,
		Timeout:   time.Duration(timeout) * time.Second,
	}
	return client.Do(request)
}

func Get(headers map[string]string, url string, timeout int) (*http.Response, error) {
	if err := checkSSRF(url); err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logc.Error(context.Background(), fmt.Sprintf("Tools get 请求建立失败, err: %s", err.Error()))
		return nil, err
	}
	for k, v := range headers {
		request.Header.Set(k, v)
	}

	resp, err := do(request, timeout)
	if err != nil {
		logc.Error(context.Background(), fmt.Sprintf("Tools get 请求发送失败, err: %s", err.Error()))
		return nil, err
	}

	return resp, nil
}

func Post(headers map[string]string, url string, bodyReader *bytes.Reader, timeout int) (*http.Response, error) {
	if err := checkSSRF(url); err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		logc.Error(context.Background(), fmt.Sprintf("Tools post 请求建立失败, err: %s", err.Error()))
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		request.Header.Set(k, v)
	}

	resp, err := do(request, timeout)
	if err != nil {
		logc.Error(context.Background(), fmt.Sprintf("Tools post 请求发送失败, err: %s", err.Error()))
		return nil, err
	}

	return resp, nil
}

// CreateBasicAuthHeader 创建带认证的HTTP头
func CreateBasicAuthHeader(username, password string) map[string]string {
	headers := make(map[string]string)
	if username != "" && password != "" {
		headers["Authorization"] = "Basic " + basicAuth(username, password)
	}
	return headers
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// MergeHeaders 合并HTTP头
func MergeHeaders(headers1, headers2 map[string]string) map[string]string {
	mergedHeaders := make(map[string]string)
	for k, v := range headers1 {
		mergedHeaders[k] = v
	}
	for k, v := range headers2 {
		mergedHeaders[k] = v
	}
	return mergedHeaders
}

// ParseReaderBody 处理请求Body
func ParseReaderBody(body io.Reader, req interface{}) error {
	newBody := body
	bodyByte, err := io.ReadAll(newBody)
	if err != nil {
		return fmt.Errorf("读取 Body 失败, err: %s", err.Error())
	}
	if err := sonic.Unmarshal(bodyByte, &req); err != nil {
		return fmt.Errorf("解析 Body 失败, body: %s, err: %s", string(bodyByte), err.Error())
	}
	return nil
}
