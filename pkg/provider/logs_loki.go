package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"
)

type LokiProvider struct {
	url            string
	timeout        int64
	ExternalLabels map[string]interface{}
}

func NewLokiClient(datasource models.AlertDataSource) (LogsFactoryProvider, error) {
	return LokiProvider{
		url:            datasource.HTTP.URL,
		timeout:        datasource.HTTP.Timeout,
		ExternalLabels: datasource.Labels,
	}, nil
}

type result struct {
	Data Data `json:"data"`
}

type Data struct {
	ResultType string   `json:"status"`
	Result     []Result `json:"result"`
}

type Result struct {
	Stream map[string]interface{} `json:"stream"`
	Values []interface{}          `json:"values"`
}

func (l LokiProvider) Query(options LogQueryOptions) ([]Logs, int, error) {
	curTime := time.Now()

	if options.Loki.Query == "" {
		return nil, 0, nil
	}

	if options.Loki.Direction == "" {
		options.Loki.Direction = "backward"
	}

	if options.Loki.Limit == 0 {
		options.Loki.Limit = 100
	}

	if options.StartAt == "" {
		duration, _ := time.ParseDuration(strconv.Itoa(1) + "h")
		options.StartAt = curTime.Add(-duration).Format(time.RFC3339Nano)
	}

	if options.EndAt == "" {
		options.EndAt = curTime.Format(time.RFC3339Nano)
	}

	args := fmt.Sprintf("/loki/api/v1/query_range?query=%s&direction=%s&limit=%d&start=%d&end=%d", url.QueryEscape(options.Loki.Query), options.Loki.Direction, options.Loki.Limit, options.StartAt.(int64), options.EndAt.(int64))
	requestURL := l.url + args
	res, err := tools.Get(nil, requestURL, 10)
	if err != nil {
		return nil, 0, err
	}

	var resultData result
	if err := tools.ParseReaderBody(res.Body, &resultData); err != nil {
		return nil, 0, errors.New(fmt.Sprintf("json.Unmarshal failed, %s", err.Error()))
	}

	var (
		count      int // count 用于统计日志条数
		data       []Logs
		streamList = []map[string]interface{}{}
		msg        []interface{}
	)
	for _, v := range resultData.Data.Result {
		streamList = append(streamList, v.Stream)
		count += len(v.Values)
		for _, m := range v.Values {
			if len(m.([]interface{})) < 2 {
				continue
			}
			msg = append(msg, m.([]interface{})[1])
		}
	}

	data = append(data, Logs{
		ProviderName: LokiDsProviderName,
		Metric:       commonKeyValuePairs(streamList),
		Message:      msg,
	})

	return data, count, nil
}

func (l LokiProvider) Check() (bool, error) {
	res, err := tools.Get(nil, l.url+"/loki/api/v1/labels", int(l.timeout))
	if err != nil {
		return false, err
	}

	if res.StatusCode != http.StatusOK {
		logc.Error(context.Background(), fmt.Errorf("unhealthy status: %d", res.StatusCode))
		return false, fmt.Errorf("unhealthy status: %d", res.StatusCode)
	}

	return true, nil
}

func (l LokiProvider) GetExternalLabels() map[string]interface{} {
	return l.ExternalLabels
}
