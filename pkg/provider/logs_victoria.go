package provider

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"watchAlert/pkg/ctx"

	"watchAlert/internal/models"
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

	if options.VictoriaLogs.Limit == 0 {
		options.VictoriaLogs.Limit = 500
	}

	args := fmt.Sprintf("/select/logsql/query?query=%s&limit=%d&start=%d&end=%d", url.QueryEscape(options.VictoriaLogs.Query), options.VictoriaLogs.Limit, options.StartAt.(int32), options.EndAt.(int32))
	requestURL := v.URL + args
	res, err := tools.Get(tools.CreateBasicAuthHeader(v.Username, v.Password), requestURL, 10)

	if err != nil {
		logc.Error(ctx.Ctx, fmt.Sprintf("查询VictoriaLogs失败: %s", err.Error()))
		return nil, 0, err
	}

	respBody, _ := io.ReadAll(res.Body)

	if res.StatusCode != 200 {
		errMsg := fmt.Sprintf("查询VictoriaLogs失败: %s", string(respBody))
		logc.Error(v.Ctx, errMsg)
		return nil, 0, fmt.Errorf(errMsg)
	}

	var entries []map[string]interface{}
	scanner := bufio.NewScanner(bytes.NewReader(respBody))
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var entry map[string]interface{}
		if err := json.Unmarshal(line, &entry); err != nil {
			return nil, 0, fmt.Errorf("解析行失败: %v，内容: %s", err, string(line))
		}
		entries = append(entries, entry)
	}

	var (
		logs []Logs
		msg  []any
	)

	for _, data := range entries {
		msg = append(msg, data["_msg"])
	}

	logs = append(logs, Logs{
		ProviderName: VictoriaLogsDsProviderName,
		Metric:       v.getMetricLabels(entries),
		Message:      msg,
	})

	return logs, len(entries), nil
}

func (v VictoriaLogsProvider) getMetricLabels(entries []map[string]interface{}) map[string]interface{} {
	metric := commonKeyValuePairs(entries)
	delete(metric, "_stream")
	delete(metric, "_stream_id")
	delete(metric, "log.file.path")
	return metric
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
