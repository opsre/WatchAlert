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
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/tools"

	"github.com/zeromicro/go-zero/core/logc"
)

type (
	VictoriaLogsProvider struct {
		URL            string         `json:"url"`
		Timeout        int64          `json:"timeout"`
		ExternalLabels map[string]any `json:"external_labels"`
		Ctx            context.Context
	}

	// VictoriaQueryResult represents the result of a query to VictoriaMetrics.
	VictoriaQueryResult struct {
		Data []struct {
			Time  int64  `json:"_time"` // 时间戳（纳秒）
			Msg   string `json:"_msg"`  // 日志消息
			Level string `json:"level"` // 日志级别
		} `json:"data"`
		Meta any `json:"meta"`
	}
)

func NewVictoriaLogsProvider(ctx context.Context, datasource models.AlertDataSource) (LogsFactoryProvider, error) {
	return VictoriaLogsProvider{
		URL:            datasource.HTTP.URL,
		Timeout:        datasource.HTTP.Timeout,
		ExternalLabels: datasource.Labels,
		Ctx:            ctx,
	}, nil
}

func (v VictoriaLogsProvider) Query(options LogQueryOptions) ([]Logs, int, error) {
	curTime := time.Now()

	if options.StartAt == "" {
		duration, _ := time.ParseDuration(strconv.Itoa(1) + "h")
		options.StartAt = curTime.Add(-duration).Format(time.RFC3339Nano)
	}

	if options.EndAt == "" {
		options.EndAt = curTime.Format(time.RFC3339Nano)
	}
	// 构造请求参数
	params := url.Values{}
	params.Add("query",
		fmt.Sprintf("%s AND _time:[%s TO %s]",
			v.URL,
			options.StartAt,
			options.EndAt,
		))

	requestURL := fmt.Sprintf("%s?%s", v.URL, params.Encode())

	res, err := tools.Get(nil, requestURL, 10)
	if err != nil {
		logc.Error(ctx.Ctx, fmt.Sprintf("查询VictoriaLogs失败: %s", err.Error()))
		return nil, 0, err
	}

	var resultData VictoriaQueryResult
	if err := tools.ParseReaderBody(res.Body, &resultData); err != nil {
		logc.Error(ctx.Ctx, fmt.Sprintf("解析VictoriaLogs结果失败: %s", err.Error()))
		return nil, 0, errors.New(fmt.Sprintf("json.Unmarshal failed, %s", err.Error()))
	}

	var logs []Logs
	var metric map[string]any
	var msg []any

	for _, data := range resultData.Data {
		metric["level"] = data.Level
		metric["_time"] = data.Time
		msg = append(msg, data.Msg)

	}
	logs = append(logs, Logs{
		ProviderName: VictoriaDsProviderName,
		Metric:       metric,
		Message:      msg,
	})
	return logs, len(resultData.Data), nil
}

func (v VictoriaLogsProvider) Check() (bool, error) {
	res, err := tools.Get(nil, v.URL+"/health", int(v.Timeout))
	if err != nil {
		return false, err
	}

	if res.StatusCode != http.StatusOK {
		logc.Error(v.Ctx, fmt.Errorf("unhealthy status: %d", res.StatusCode))
		return false, fmt.Errorf("unhealthy status: %d", res.StatusCode)
	}
	return true, nil
}

func (v VictoriaLogsProvider) GetExternalLabels() map[string]interface{} {
	return v.ExternalLabels
}
