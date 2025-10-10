package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
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

func (l LokiProvider) Query(options LogQueryOptions) (Logs, int, error) {
	curTime := time.Now()

	if options.Loki.Query == "" {
		return Logs{}, 0, nil
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
		return Logs{}, 0, err
	}

	var resultData result
	if err := tools.ParseReaderBody(res.Body, &resultData); err != nil {
		return Logs{}, 0, errors.New(fmt.Sprintf("json.Unmarshal failed, %s", err.Error()))
	}

	var (
		count   int // count 用于统计日志条数
		message []map[string]interface{}
	)
	for _, v := range resultData.Data.Result {
		count += len(v.Values)
		/*
				"values": [
			          [
			            "1746671236062002147",
			            //"{\"level\":\"INFO\",\"time\":\"2025-05-08T02:27:16.061Z\",\"pid\":1,\"hostname\":\"hedwig-5f8fcc9c68-wgplm\",\"req\":{\"id\":2925131,\"method\":\"GET\",\"url\":\"/api/health/check\",\"query\":{},\"params\":{\"0\":\"api/health/check\"},\"headers\":{\"host\":\"10.42.0.179:8080\",\"user-agent\":\"kube-probe/1.20\",
					  ],
				]
		*/

		if len(v.Values) == 0 {
			continue
		}

		for _, m := range v.Values {
			firstValue, ok := m.([]interface{})
			if !ok || len(firstValue) < 2 {
				logc.Error(context.Background(), "Loki - Values[0] 类型错误或长度不足")
				continue
			}

			rawLog := firstValue[1]
			var jsonData []byte

			switch val := rawLog.(type) {
			case []byte:
				jsonData = val
			case string:
				jsonData = []byte(val)
			default:
				logc.Error(context.Background(), "Loki - Values[0][1] 类型不是 []byte 或 string")
				continue
			}

			var msg map[string]interface{}
			err := sonic.Unmarshal(jsonData, &msg)
			if err != nil {
				logc.Error(context.Background(), fmt.Sprintf("解析 Loki 日志数据错误, %v", string(jsonData)))
				continue
			}
			message = append(message, msg)
		}
	}

	return Logs{
		ProviderName: LokiDsProviderName,
		Message:      message,
	}, count, nil
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
