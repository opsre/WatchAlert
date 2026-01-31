package provider

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"

	"github.com/zeromicro/go-zero/core/logc"
)

type PrometheusProvider struct {
	client         v1.API
	ExternalLabels map[string]interface{}
	Address        string
	Username       string
	Password       string
	Headers        map[string]string
	Timeout        int64
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
		ExternalLabels: ds.Labels,
		Username:       ds.Auth.User,
		Password:       ds.Auth.Pass,
		Headers:        ds.HTTP.Headers,
		Timeout:        ds.HTTP.Timeout,
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
			Metric:    metric,
			Value:     float64(item.Value),
			Timestamp: float64(item.Timestamp),
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
				Timestamp: float64(value.Timestamp),
				Value:     float64(value.Value),
				Metric:    metric,
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
