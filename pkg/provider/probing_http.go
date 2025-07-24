package provider

import (
	"bytes"
	"net/http"
	"time"
	"watchAlert/pkg/tools"
)

const (
	GetHTTPMethod  = "GET"
	PostHTTPMethod = "POST"
)

type HTTPer struct{}

func NewEndpointHTTPer() EndpointFactoryProvider {
	return HTTPer{}
}

func (h HTTPer) Pilot(option EndpointOption) (EndpointValue, error) {
	var (
		ev      EndpointValue
		res     *http.Response
		err     error
		headers = make(map[string][]string)
	)

	// 开始时间
	start := time.Now()
	switch option.HTTP.Method {
	case GetHTTPMethod:
		res, err = tools.Get(option.HTTP.Header, option.Endpoint, option.Timeout)
		if err != nil {
			return ev, err
		}
		headers = res.Header
		defer res.Body.Close()
	case PostHTTPMethod:
		res, err = tools.Post(option.HTTP.Header, option.Endpoint, bytes.NewReader([]byte(option.HTTP.Body)), option.Timeout)
		if err != nil {
			return ev, err
		}
		headers = res.Header
		defer res.Body.Close()
	}
	end := time.Now()
	// 计算请求耗时
	latency := end.Sub(start).Milliseconds()

	return convertHTTPerToEndpointValue(HttperInformation{
		Address:    res.Request.URL.String(),
		StatusCode: float64(res.StatusCode),
		Latency:    float64(latency),
		Headers:    headers,
	}), nil
}

func convertHTTPerToEndpointValue(detail HttperInformation) EndpointValue {
	headers := make(map[string]interface{})
	for k, v := range detail.Headers {
		headers[k] = v[0]
	}

	return EndpointValue{
		"address":    detail.Address,
		"StatusCode": detail.StatusCode,
		"Latency":    detail.Latency,
		"headers":    headers,
	}
}
