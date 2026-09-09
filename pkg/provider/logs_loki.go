package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"

	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logc"
)

type LokiProvider struct {
	Url            string
	Timeout        int64
	ExternalLabels map[string]interface{}
	Headers        map[string]string
}

func NewLokiClient(datasource models.AlertDataSource) (LogsFactoryProvider, error) {
	return LokiProvider{
		Url:            datasource.HTTP.URL,
		Timeout:        datasource.HTTP.Timeout,
		ExternalLabels: datasource.Labels,
		Headers:        datasource.HTTP.Headers,
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

// parseLokiTimeParam 将 LogQueryOptions.StartAt/EndAt 转换为 Unix 时间戳（秒）。
// 支持 int64/int32/int（Unix 秒）以及字符串（Unix 秒数字符串或 RFC3339Nano 格式）。
// 当传入值为空或无法解析时，返回 (0, false)。
func parseLokiTimeParam(v interface{}) (int64, bool) {
	switch val := v.(type) {
	case int64:
		if val == 0 {
			return 0, false
		}
		return val, true
	case int32:
		if val == 0 {
			return 0, false
		}
		return int64(val), true
	case int:
		if val == 0 {
			return 0, false
		}
		return int64(val), true
	case string:
		if val == "" {
			return 0, false
		}
		// 尝试解析为 Unix 时间戳字符串
		if unix, err := strconv.ParseInt(val, 10, 64); err == nil {
			return unix, true
		}
		// 尝试解析为 RFC3339Nano 格式
		if t, err := time.Parse(time.RFC3339Nano, val); err == nil {
			return t.Unix(), true
		}
		// 尝试解析为 RFC3339 格式
		if t, err := time.Parse(time.RFC3339, val); err == nil {
			return t.Unix(), true
		}
		return 0, false
	default:
		return 0, false
	}
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

	startAt, ok := parseLokiTimeParam(options.StartAt)
	if !ok {
		duration, _ := time.ParseDuration(strconv.Itoa(1) + "h")
		startAt = curTime.Add(-duration).Unix()
	}

	endAt, ok := parseLokiTimeParam(options.EndAt)
	if !ok {
		endAt = curTime.Unix()
	}

	args := fmt.Sprintf("/loki/api/v1/query_range?query=%s&direction=%s&limit=%d&start=%d&end=%d", url.QueryEscape(options.Loki.Query), options.Loki.Direction, options.Loki.Limit, startAt, endAt)
	requestURL := l.Url + args

	var headers = make(map[string]string)
	for key, value := range l.Headers {
		headers[key] = value
	}

	res, err := tools.Get(headers, requestURL, 10)
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
	var headers = make(map[string]string)
	for key, value := range l.Headers {
		headers[key] = value
	}

	res, err := tools.Get(headers, l.Url+"/loki/api/v1/labels", int(l.Timeout))
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
