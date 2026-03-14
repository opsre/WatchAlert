package provider

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"

	"github.com/zeromicro/go-zero/core/logc"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/prometheus/prompb"
)

type PrometheusProvider struct {
	client         v1.API
	ExternalLabels map[string]interface{}
	Address        string
	WriteURL       string
	Username       string
	Password       string
	Headers        map[string]string
	Timeout        int64
	httpClient     *http.Client
}

// authenticatedTransport 包装 http.RoundTripper 以添加认证头和额外的headers
type authenticatedTransport struct {
	Transport http.RoundTripper
	Username  string
	Password  string
	Headers   map[string]string
}

// RoundTrip 实现 http.RoundTripper 接口
func (t *authenticatedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.Username != "" && t.Password != "" {
		req.SetBasicAuth(t.Username, t.Password)
	}

	for key, value := range t.Headers {
		req.Header.Set(key, value)
	}

	return t.Transport.RoundTrip(req)
}

func NewPrometheusClient(ds models.AlertDataSource) (MetricsFactoryProvider, error) {
	transport := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}

	var roundTripper http.RoundTripper = transport
	if ds.Auth.User != "" || ds.Auth.Pass != "" || len(ds.HTTP.Headers) > 0 {
		roundTripper = &authenticatedTransport{
			Transport: transport,
			Username:  ds.Auth.User,
			Password:  ds.Auth.Pass,
			Headers:   ds.HTTP.Headers,
		}
	}

	clientConfig := api.Config{
		Address:      ds.HTTP.URL,
		RoundTripper: roundTripper,
	}

	client, err := api.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}

	return PrometheusProvider{
		client:         v1.NewAPI(client),
		Address:        ds.HTTP.URL,
		WriteURL:       ds.Write.URL,
		ExternalLabels: ds.Labels,
		Username:       ds.Auth.User,
		Password:       ds.Auth.Pass,
		Headers:        ds.HTTP.Headers,
		Timeout:        ds.HTTP.Timeout,
		httpClient:     &http.Client{Transport: roundTripper, Timeout: time.Duration(ds.HTTP.Timeout) * time.Second},
	}, nil
}

type QueryResponse struct {
	Status     string     `json:"status"`
	MetricData MetricData `json:"data"`
}

type MetricData struct {
	MetricResult []MetricResult `json:"result"`
	ResultType   string         `json:"resultType"`
}

type MetricResult struct {
	Metric map[string]interface{} `json:"metric"`
	Value  []interface{}          `json:"value"`
	Values [][]interface{}        `json:"values"`
}

func (v PrometheusProvider) Query(promQL string) ([]Metrics, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(v.Timeout)*time.Second)
	defer cancel()
	result, _, err := v.client.Query(ctx, promQL, time.Now(), v1.WithTimeout(time.Duration(v.Timeout)*time.Second))
	if err != nil {
		return nil, err
	}
	return Vectors(result), nil
}

func (v PrometheusProvider) QueryRange(promQL string, start, end time.Time, step time.Duration) ([]Metrics, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(v.Timeout)*time.Second)
	defer cancel()

	r := v1.Range{
		Start: start,
		End:   end,
		Step:  step,
	}

	result, _, err := v.client.QueryRange(ctx, promQL, r, v1.WithTimeout(time.Duration(v.Timeout)*time.Second))
	if err != nil {
		return nil, err
	}

	return Matrix(result), nil
}

func Vectors(value model.Value) []Metrics {
	var vectors []Metrics
	items, ok := value.(model.Vector)
	if !ok {
		return []Metrics{}
	}

	for _, item := range items {
		if math.IsNaN(float64(item.Value)) {
			logc.Infof(context.Background(), "Skipping NaN or Inf value: %v", item.Value)
			continue
		}

		var metric = make(map[string]interface{})
		for k, v := range item.Metric {
			metric[string(k)] = string(v)
		}

		vectors = append(vectors, Metrics{
			Labels:    metric,
			Value:     float64(item.Value),
			Timestamp: int64(item.Timestamp),
		})
	}

	return vectors
}

// Matrix 将 Prometheus QueryRange 结果转换为 Metrics 列表
func Matrix(value model.Value) []Metrics {
	var metrics []Metrics
	matrix, ok := value.(model.Matrix)
	if !ok {
		return []Metrics{}
	}

	for _, stream := range matrix {
		var metric = make(map[string]interface{})
		for k, v := range stream.Metric {
			metric[string(k)] = string(v)
		}

		for _, value := range stream.Values {
			if math.IsNaN(float64(value.Value)) {
				continue
			}

			metrics = append(metrics, Metrics{
				Timestamp: int64(value.Timestamp),
				Value:     float64(value.Value),
				Labels:    metric,
			})
		}
	}

	return metrics
}

func (v PrometheusProvider) Check() (bool, error) {
	var headers map[string]string
	checkURL := v.Address + "/api/v1/query?query=1%2B1"
	if v.Username != "" && v.Password != "" {
		headers = tools.CreateBasicAuthHeader(v.Username, v.Password)
	}
	headers = tools.MergeHeaders(headers, v.Headers)
	res, err := tools.Get(headers, checkURL, int(v.Timeout))
	if err != nil {
		logc.Errorf(context.Background(), "Health check failed, URL: %s, Error: %v", checkURL, err)
		return false, fmt.Errorf("health check failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		logc.Errorf(context.Background(), "Health check received unhealthy status: %d, URL: %s", res.StatusCode, checkURL)
		return false, fmt.Errorf("unhealthy status: %d", res.StatusCode)
	}
	return true, nil
}

func (v PrometheusProvider) GetExternalLabels() map[string]interface{} {
	return v.ExternalLabels
}

// Write 将记录规则结果写入 Prometheus 远程写入端点
func (v PrometheusProvider) Write(ctx context.Context, metrics []Metrics, externalLabels map[string]string) error {
	if len(metrics) == 0 {
		return nil
	}

	// 转换为Prometheus Remote Write格式
	writeRequest, err := v.convertToRemoteWriteFormat(metrics, externalLabels)
	if err != nil {
		return fmt.Errorf("转换为Remote Write格式失败: %w", err)
	}

	// 序列化为Protobuf
	data, err := proto.Marshal(writeRequest)
	if err != nil {
		return fmt.Errorf("序列化Protobuf失败: %w", err)
	}

	// Snappy压缩
	compressed := snappy.Encode(nil, data)

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, v.WriteURL, bytes.NewBuffer(compressed))
	if err != nil {
		return fmt.Errorf("创建Prometheus请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/x-protobuf")
	req.Header.Set("Content-Encoding", "snappy")
	req.Header.Set("X-Prometheus-Remote-Write-Version", "0.1.0")

	// 设置认证
	if v.Username != "" && v.Password != "" {
		req.SetBasicAuth(v.Username, v.Password)
	}

	// 发送请求
	resp, err := v.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("发送Prometheus请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	// 处理错误响应
	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("prometheus写入失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
}

// convertToRemoteWriteFormat 转换为 Remote Write格式
func (v PrometheusProvider) convertToRemoteWriteFormat(results []Metrics, externalLabels map[string]string) (*prompb.WriteRequest, error) {
	var timeSeries []prompb.TimeSeries

	for _, metric := range results {
		// 构建标签
		labels := []prompb.Label{
			{Name: "__name__", Value: metric.Name},
		}

		// 添加其他标签
		for k, v := range metric.Labels {
			labels = append(labels, prompb.Label{
				Name:  k,
				Value: fmt.Sprintf("%v", v),
			})
		}

		for k, v := range externalLabels {
			labels = append(labels, prompb.Label{
				Name:  k,
				Value: fmt.Sprintf("%v", v),
			})
		}

		// 创建样本
		sample := prompb.Sample{
			Value:     metric.Value,
			Timestamp: int64(time.Now().UnixMilli()),
		}

		// 创建时间序列
		ts := prompb.TimeSeries{
			Labels:  labels,
			Samples: []prompb.Sample{sample},
		}

		timeSeries = append(timeSeries, ts)
	}

	// 构造WriteRequest
	writeReq := &prompb.WriteRequest{
		Timeseries: timeSeries,
	}

	return writeReq, nil
}
