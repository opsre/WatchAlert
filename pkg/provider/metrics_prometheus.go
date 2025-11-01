package provider

import (
	"context"
	"math"
	"net/http"
	"time"
	"watchAlert/internal/models"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type PrometheusProvider struct {
	ExternalLabels map[string]interface{}
	apiV1          v1.API
}

// BasicAuthTransport 实现带认证的HTTP传输层
type BasicAuthTransport struct {
	Username string
	Password string
	Base     http.RoundTripper
}

func (t *BasicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.Username != "" || t.Password != "" {
		req.SetBasicAuth(t.Username, t.Password)
	}
	return t.Base.RoundTrip(req)
}

func NewPrometheusClient(source models.AlertDataSource) (MetricsFactoryProvider, error) {
	// 创建基础传输层
	baseTransport := http.DefaultTransport

	// 配置认证传输层
	authTransport := &BasicAuthTransport{
		Username: source.Auth.User,
		Password: source.Auth.Pass,
		Base:     baseTransport,
	}

	// 创建客户端配置
	clientConfig := api.Config{
		Address:      source.HTTP.URL,
		RoundTripper: authTransport,
	}

	// 创建带认证的客户端
	client, err := api.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}

	return PrometheusProvider{
		apiV1:          v1.NewAPI(client),
		ExternalLabels: source.Labels,
	}, nil
}

func (p PrometheusProvider) Query(promQL string) ([]Metrics, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, _, err := p.apiV1.Query(ctx, promQL, time.Now(), v1.WithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}

	return ConvertVectors(result), nil
}

func (p PrometheusProvider) QueryRange(promQL string, start, end time.Time, step time.Duration) ([]Metrics, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	r := v1.Range{
		Start: start,
		End:   end,
		Step:  step,
	}

	result, _, err := p.apiV1.QueryRange(ctx, promQL, r, v1.WithTimeout(20*time.Second))
	if err != nil {
		return nil, err
	}

	return ConvertMatrix(result), nil
}

func ConvertVectors(value model.Value) (lst []Metrics) {
	items, ok := value.(model.Vector)
	if !ok {
		return
	}

	for _, item := range items {
		if math.IsNaN(float64(item.Value)) {
			continue
		}

		var metric = make(map[string]interface{})
		for k, v := range item.Metric {
			metric[string(k)] = string(v)
		}

		lst = append(lst, Metrics{
			Timestamp: float64(item.Timestamp),
			Value:     float64(item.Value),
			Metric:    metric,
		})
	}
	return
}

// ConvertMatrix 将 Prometheus QueryRange 结果转换为 Metrics 列表
func ConvertMatrix(value model.Value) (lst []Metrics) {
	matrix, ok := value.(model.Matrix)
	if !ok {
		return
	}

	for _, stream := range matrix {
		var metric = make(map[string]interface{})
		for k, v := range stream.Metric {
			metric[string(k)] = string(v)
		}

		// 将每个时间点的数据转换为单独的 Metrics
		for _, value := range stream.Values {
			if math.IsNaN(float64(value.Value)) {
				continue
			}

			lst = append(lst, Metrics{
				Timestamp: float64(value.Timestamp),
				Value:     float64(value.Value),
				Metric:    metric,
			})
		}
	}
	return
}

func (p PrometheusProvider) Check() (bool, error) {
	_, err := p.apiV1.Config(context.Background())
	if err != nil {
		return false, err
	}

	return true, nil
}

func (p PrometheusProvider) GetExternalLabels() map[string]interface{} {
	return p.ExternalLabels
}
