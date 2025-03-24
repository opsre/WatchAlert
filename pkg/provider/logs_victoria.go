package provider

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
		Username       string `json:"username"`
		Password       string `json:"password"`
	}

	// VictoriaQueryResult represents the result of a query to VictoriaMetrics.
	VictoriaQueryResult struct {
		Time      time.Time `json:"_time"` // 时间戳（纳秒）
		Msg       string    `json:"_msg"`  // 日志消息
		StreamId  string    `json:"_stream_id"`
		AgentName string    `json:"agent.name"`
		AgentType string    `json:"agent.type"`
	}
)

// NewVictoriaLogsClient 创建一个新的 VictoriaLogsProvider 实例。
func NewVictoriaLogsClient(ds models.AlertDataSource) (VictoriaLogsProvider, error) {
	return VictoriaLogsProvider{
		URL:            ds.HTTP.URL,
		ExternalLabels: ds.Labels,
		Username:       ds.Auth.User,
		Password:       ds.Auth.Pass,
	}, nil
}

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

	if options.Victoria.Limit == "" {
		options.Victoria.Limit = "50"
	}

	// 构造请求参数
	params := make(map[string]any)
	params["query"] = "*"
	params["limit"] = 50

	params["start"] = options.StartAt.(int32) // 开始时间
	params["end"] = options.EndAt.(int32)     // 结束时间

	args := fmt.Sprintf("/select/logsql/query?query=%s&limit=%s&start=%d&end=%d", options.Victoria.Query, options.Victoria.Limit, options.StartAt, options.EndAt)

	requestURL := v.URL + args
	res, err := tools.Get(nil, requestURL, 10)

	respBody, _ := io.ReadAll(res.Body)

	if err != nil {
		logc.Error(ctx.Ctx, fmt.Sprintf("查询VictoriaLogs失败: %s", err.Error()))
		return nil, 0, err
	}

	var entries []VictoriaQueryResult
	scanner := bufio.NewScanner(bytes.NewReader(respBody))
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var entry VictoriaQueryResult
		if err := json.Unmarshal(line, &entry); err != nil {
			return nil, 0, fmt.Errorf("解析行失败: %v，内容: %s", err, string(line))
		}
		entries = append(entries, entry)
	}

	var logs []Logs
	metric := make(map[string]any, len(entries))
	var msg []any
	for _, data := range entries {
		metric["_stream_id"] = data.StreamId
		metric["agent.name"] = data.AgentName
		metric["agent.type"] = data.AgentType
		metric["_time"] = data.Time
		msg = append(msg, data.Msg)

	}

	logs = append(logs, Logs{
		ProviderName: VictoriaLogsDsProviderName,
		Metric:       metric,
		Message:      msg,
	})

	return logs, len(entries), nil
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
