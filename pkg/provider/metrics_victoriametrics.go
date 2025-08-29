package provider

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"watchAlert/internal/models"
	utilsHttp "watchAlert/pkg/tools"
)

type VictoriaMetricsProvider struct {
	ExternalLabels map[string]interface{}
	address        string
	username       string
	password       string
}

func NewVictoriaMetricsClient(ds models.AlertDataSource) (MetricsFactoryProvider, error) {
	return VictoriaMetricsProvider{
		address:        ds.HTTP.URL,
		ExternalLabels: ds.Labels,
		username:       ds.Auth.User,
		password:       ds.Auth.Pass,
	}, nil
}

type QueryResponse struct {
	Status string `json:"status"`
	VMData VMData `json:"data"`
}

type VMData struct {
	VMResult   []VMResult `json:"result"`
	ResultType string     `json:"resultType"`
}

type VMResult struct {
	Metric map[string]interface{} `json:"metric"`
	Value  []interface{}          `json:"value"`
}

func (v VictoriaMetricsProvider) Query(promQL string) ([]Metrics, error) {
	params := url.Values{}
	params.Add("query", promQL)
	params.Add("time", strconv.FormatInt(time.Now().Unix(), 10))
	fullURL := fmt.Sprintf("%s%s?%s", v.address, "/api/v1/query", params.Encode())

	// 创建带认证的HTTP请求
	resp, err := utilsHttp.Get(utilsHttp.CreateBasicAuthHeader(v.username, v.password), fullURL, 10)
	if err != nil {
		logc.Error(context.Background(), "VictoriaMetrics query failed", "error", err)
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var vmRespBody QueryResponse
	if err := utilsHttp.ParseReaderBody(resp.Body, &vmRespBody); err != nil {
		logc.Error(context.Background(), "Parse response failed", "error", err)
		return nil, fmt.Errorf("parse response failed: %w", err)
	}

	return vmVectors(vmRespBody.VMData.VMResult), nil
}

func vmVectors(res []VMResult) []Metrics {
	var vectors []Metrics
	for _, item := range res {
		if len(item.Value) < 2 {
			continue
		}

		timestamp, ok1 := item.Value[0].(float64)
		valueStr, ok2 := item.Value[1].(string)
		if !ok1 || !ok2 {
			logc.Error(context.Background(), "Invalid value format")
			continue
		}

		valueFloat, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			logc.Error(context.Background(), "Value conversion failed", "error", err)
			continue
		}

		vectors = append(vectors, Metrics{
			Metric:    item.Metric,
			Value:     valueFloat,
			Timestamp: timestamp,
		})
	}
	return vectors
}

func (v VictoriaMetricsProvider) Check() (bool, error) {
	res, err := utilsHttp.Get(utilsHttp.CreateBasicAuthHeader(v.username, v.password), v.address+"/api/v1/query?query=1%2B1", 10)
	if err != nil {
		logc.Error(context.Background(), fmt.Errorf("health check failed: %w", err))
		return false, fmt.Errorf("health check failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		logc.Error(context.Background(), fmt.Errorf("unhealthy status: %d", res.StatusCode))
		return false, fmt.Errorf("unhealthy status: %d", res.StatusCode)
	}
	return true, nil
}

func (v VictoriaMetricsProvider) GetExternalLabels() map[string]interface{} {
	return v.ExternalLabels
}
