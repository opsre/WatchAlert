package provider

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"io"
	"net/http"
	"net/url"
	"time"
	"watchAlert/internal/ctx"
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
func NewVictoriaLogsClient(ctx context.Context, datasource models.AlertDataSource) (LogsFactoryProvider, error) {
	return VictoriaLogsProvider{
		URL:            datasource.HTTP.URL,
		Timeout:        datasource.HTTP.Timeout,
		ExternalLabels: datasource.Labels,
		Username:       datasource.Auth.User,
		Password:       datasource.Auth.Pass,
		Ctx:            ctx,
	}, nil
}

func (v VictoriaLogsProvider) Query(options LogQueryOptions) (Logs, int, error) {
	curTime := time.Now()

	if options.StartAt == "" || options.StartAt == nil {
		options.StartAt = int32(tools.ParserDuration(curTime, 30, "m").Unix())
	}

	if options.EndAt == "" || options.EndAt == nil {
		options.EndAt = int32(curTime.Unix())
	}

	if options.VictoriaLogs.Limit == 0 {
		options.VictoriaLogs.Limit = 500
	}

	args := fmt.Sprintf("/select/logsql/query?query=%s&limit=%d&start=%d&end=%d", url.QueryEscape(options.VictoriaLogs.Query), options.VictoriaLogs.Limit, options.StartAt.(int32), options.EndAt.(int32))
	requestURL := v.URL + args
	res, err := tools.Get(tools.CreateBasicAuthHeader(v.Username, v.Password), requestURL, 10)

	if err != nil {
		logc.Error(ctx.Ctx, fmt.Sprintf("查询VictoriaLogs失败: %s", err.Error()))
		return Logs{}, 0, err
	}

	respBody, _ := io.ReadAll(res.Body)

	if res.StatusCode != 200 {
		errMsg := fmt.Sprintf("查询VictoriaLogs失败: %s", string(respBody))
		logc.Error(v.Ctx, errMsg)
		return Logs{}, 0, fmt.Errorf(errMsg)
	}

	var (
		message []map[string]interface{}
		count   int
	)
	scanner := bufio.NewScanner(bytes.NewReader(respBody))
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var msg map[string]interface{}
		if err := sonic.Unmarshal(line, &msg); err != nil {
			logc.Error(context.Background(), fmt.Sprintf("VictoriaLogs - 解析行失败: %v，内容: %s", err, string(line)))
			continue
		}
		message = append(message, msg)
		count++
	}

	return Logs{
		ProviderName: VictoriaLogsDsProviderName,
		Message:      message,
	}, count, nil
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
