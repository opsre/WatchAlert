package provider

import (
	"fmt"
	"io"
	"strconv"
	"time"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"
)

type JaegerDsProvider struct {
	ExternalLabels map[string]interface{}
	url            string
	headers        map[string]string
}

func NewJaegerClient(datasource models.AlertDataSource) (TracesFactoryProvider, error) {
	headers := tools.MergeHeaders(tools.CreateBasicAuthHeader(datasource.Auth.User, datasource.Auth.Pass), datasource.HTTP.Headers)
	res, err := tools.Get(headers, datasource.HTTP.URL, 10)
	if err != nil {
		return JaegerDsProvider{}, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return JaegerDsProvider{}, fmt.Errorf("Jaeger health check failed, status: %d", res.StatusCode)
	}

	return JaegerDsProvider{
		url:            datasource.HTTP.URL,
		headers:        headers,
		ExternalLabels: datasource.Labels,
	}, nil
}

type JaegerResult struct {
	Data []JaegerData `json:"data"`
}

type JaegerData struct {
	TraceId string `json:"traceID"`
}

func (j JaegerDsProvider) Query(options TraceQueryOptions) ([]Traces, error) {
	curTime := time.Now()

	if options.Limit == 0 {
		options.Limit = 100
	}

	if options.StartAt == 0 {
		duration, _ := time.ParseDuration(strconv.Itoa(1) + "h")
		options.StartAt = curTime.Add(-duration).UnixMicro()
	}

	if options.EndAt == 0 {
		options.EndAt = curTime.UnixMicro()
	}

	args := fmt.Sprintf("/api/traces?service=%s&start=%d&end=%d&limit=%d&tags=%s", options.Service, options.StartAt, options.EndAt, options.Limit, options.Tags)
	requestURL := j.url + args
	res, err := tools.Get(j.headers, requestURL, 10)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("Jaeger query failed, status: %d, body: %s", res.StatusCode, string(body))
	}

	var jaegerResult JaegerResult
	if err := tools.ParseReaderBody(res.Body, &jaegerResult); err != nil {
		return nil, err
	}

	var data []Traces
	for _, t := range jaegerResult.Data {
		data = append(data, Traces{
			Service: options.Service,
			TraceId: t.TraceId,
		})
	}

	return data, nil
}

func (j JaegerDsProvider) Check() (bool, error) {
	res, err := tools.Get(j.headers, j.url, 10)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return false, fmt.Errorf("Jaeger health check failed, status: %d", res.StatusCode)
	}
	return true, nil
}

func (j JaegerDsProvider) GetServiceNames() (TraceServiceData, error) {
	url := j.url + "/api/services"
	res, err := tools.Get(j.headers, url, 10)
	if err != nil {
		return TraceServiceData{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		b, _ := io.ReadAll(res.Body)
		return TraceServiceData{}, fmt.Errorf("后端服务请求异常, Status: %d, Msg: %s", res.StatusCode, string(b))
	}

	var resData TraceServiceData
	if err := tools.ParseReaderBody(res.Body, &resData); err != nil {
		return TraceServiceData{}, fmt.Errorf("json.Unmarshal failed, %s", err.Error())
	}

	return resData, nil
}

// GetJaegerService is kept for callers using the concrete provider.
func (j JaegerDsProvider) GetJaegerService() (JaegerServiceData, error) {
	return j.GetServiceNames()
}

func (j JaegerDsProvider) GetExternalLabels() map[string]interface{} {
	return j.ExternalLabels
}
