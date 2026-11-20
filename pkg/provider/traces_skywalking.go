package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"
)

const skyWalkingGraphQLPath = "/graphql"

type SkyWalkingDsProvider struct {
	ExternalLabels map[string]interface{}
	url            string
	headers        map[string]string
}

func NewSkyWalkingClient(datasource models.AlertDataSource) (TracesFactoryProvider, error) {
	baseURL := strings.TrimRight(datasource.HTTP.URL, "/")
	headers := tools.MergeHeaders(tools.CreateBasicAuthHeader(datasource.Auth.User, datasource.Auth.Pass), datasource.HTTP.Headers)
	provider := SkyWalkingDsProvider{
		url:            baseURL,
		headers:        headers,
		ExternalLabels: datasource.Labels,
	}
	if ok, err := provider.Check(); err != nil || !ok {
		if err == nil {
			err = fmt.Errorf("SkyWalking health check failed")
		}
		return SkyWalkingDsProvider{}, err
	}
	return provider, nil
}

type skyWalkingGraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type skyWalkingGraphQLError struct {
	Message string `json:"message"`
}

type skyWalkingGraphQLResponse struct {
	Data   json.RawMessage          `json:"data"`
	Errors []skyWalkingGraphQLError `json:"errors"`
}

type skyWalkingService struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type skyWalkingServicesResponse struct {
	GetAllServices []skyWalkingService `json:"getAllServices"`
}

type skyWalkingTraceResponse struct {
	QueryTraces struct {
		Traces []struct {
			Spans []struct {
				TraceID string `json:"traceId"`
			} `json:"spans"`
		} `json:"traces"`
	} `json:"queryTraces"`
}

func (s SkyWalkingDsProvider) graphql(request skyWalkingGraphQLRequest, target interface{}) error {
	body, err := json.Marshal(request)
	if err != nil {
		return err
	}
	res, err := tools.Post(s.headers, s.url+skyWalkingGraphQLPath, bytes.NewReader(body), 10)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("SkyWalking request failed, status: %d, body: %s", res.StatusCode, string(responseBody))
	}

	var response skyWalkingGraphQLResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return fmt.Errorf("解析 SkyWalking 响应失败: %w", err)
	}
	if len(response.Errors) > 0 {
		return fmt.Errorf("SkyWalking GraphQL query failed: %s", response.Errors[0].Message)
	}
	if target == nil {
		return nil
	}
	if len(response.Data) == 0 || string(response.Data) == "null" {
		return fmt.Errorf("SkyWalking 响应缺少 data")
	}
	return json.Unmarshal(response.Data, target)
}

func (s SkyWalkingDsProvider) Check() (bool, error) {
	err := s.graphql(skyWalkingGraphQLRequest{
		Query: "query { __typename }",
	}, nil)
	return err == nil, err
}

func (s SkyWalkingDsProvider) queryServices(start, end time.Time) ([]skyWalkingService, error) {
	var response skyWalkingServicesResponse
	err := s.graphql(skyWalkingGraphQLRequest{
		Query: `query GetAllServices($duration: Duration!) {
			getAllServices(duration: $duration) { id name }
		}`,
		Variables: map[string]interface{}{
			"duration": skyWalkingDuration(start, end),
		},
	}, &response)
	if err != nil {
		return nil, err
	}
	return response.GetAllServices, nil
}

func (s SkyWalkingDsProvider) GetServiceNames() (TraceServiceData, error) {
	end := time.Now()
	start := end.Add(-24 * time.Hour)
	serviceList, err := s.queryServices(start, end)
	if err != nil {
		return TraceServiceData{}, err
	}

	services := make([]string, 0, len(serviceList))
	for _, service := range serviceList {
		if service.Name != "" {
			services = append(services, service.Name)
		}
	}
	return TraceServiceData{Data: services}, nil
}

func (s SkyWalkingDsProvider) getServiceID(serviceName string, start, end time.Time) (string, error) {
	serviceList, err := s.queryServices(start, end)
	if err != nil {
		return "", err
	}
	for _, service := range serviceList {
		if service.Name == serviceName {
			return service.ID, nil
		}
	}
	return "", fmt.Errorf("SkyWalking 服务不存在: %s", serviceName)
}

func (s SkyWalkingDsProvider) Query(options TraceQueryOptions) ([]Traces, error) {
	end := time.Now()
	start := end.Add(-time.Hour)
	if options.StartAt != 0 {
		start = time.UnixMicro(options.StartAt)
	}
	if options.EndAt != 0 {
		end = time.UnixMicro(options.EndAt)
	}
	if options.Limit == 0 {
		options.Limit = 100
	}

	tags, err := parseSkyWalkingTags(options.Tags)
	if err != nil {
		return nil, err
	}
	traceState := "ERROR"
	if len(tags) > 0 {
		traceState = "ALL"
	}

	serviceID, err := s.getServiceID(options.Service, start, end)
	if err != nil {
		return nil, err
	}

	var response skyWalkingTraceResponse
	err = s.graphql(skyWalkingGraphQLRequest{
		Query: `query queryV2Traces($condition: TraceQueryCondition) {
			queryTraces(condition: $condition) {
				traces { spans { traceId } }
			}
		}`,
		Variables: map[string]interface{}{
			"condition": map[string]interface{}{
				"serviceId":     serviceID,
				"queryDuration": skyWalkingDuration(start, end),
				"traceState":    traceState,
				"queryOrder":    "BY_START_TIME",
				"tags":          tags,
				"paging": map[string]int{
					"pageNum":  1,
					"pageSize": int(options.Limit),
				},
			},
		},
	}, &response)
	if err != nil {
		return nil, err
	}

	traces := make([]Traces, 0)
	seen := make(map[string]struct{})
	for _, item := range response.QueryTraces.Traces {
		for _, span := range item.Spans {
			traceID := span.TraceID
			if traceID == "" {
				continue
			}
			if _, ok := seen[traceID]; ok {
				continue
			}
			seen[traceID] = struct{}{}
			traces = append(traces, Traces{Service: options.Service, TraceId: traceID})
		}
	}
	return traces, nil
}

func skyWalkingDuration(start, end time.Time) map[string]interface{} {
	return map[string]interface{}{
		"start":     start.Format("2006-01-02 150405"),
		"end":       end.Format("2006-01-02 150405"),
		"step":      "SECOND",
		"coldStage": false,
	}
}

func parseSkyWalkingTags(raw string) ([]map[string]string, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}

	var values map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &values); err != nil {
		return nil, fmt.Errorf("SkyWalking 标签条件必须是 JSON 对象: %w", err)
	}
	tags := make([]map[string]string, 0, len(values))
	for key, value := range values {
		if key == "" {
			continue
		}
		var stringValue string
		switch value := value.(type) {
		case string:
			stringValue = value
		case float64:
			stringValue = strconv.FormatFloat(value, 'f', -1, 64)
		case bool:
			stringValue = strconv.FormatBool(value)
		default:
			return nil, fmt.Errorf("SkyWalking 标签值必须是字符串、数字或布尔值: %s", key)
		}
		tags = append(tags, map[string]string{"key": key, "value": stringValue})
	}
	return tags, nil
}

func (s SkyWalkingDsProvider) GetExternalLabels() map[string]interface{} {
	return s.ExternalLabels
}
