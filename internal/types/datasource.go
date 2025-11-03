package types

import (
	"fmt"
	"strings"
	"time"
	"watchAlert/internal/models"
)

const (
	ErrorQueryIsEmpty = "query is empty"
)

type RequestDatasourceCreate struct {
	TenantId         string                    `json:"tenantId"`
	Name             string                    `json:"name"`
	Labels           map[string]interface{}    `json:"labels"` // 额外标签，会添加到事件Metric中，可用于区分数据来源；
	Type             string                    `json:"type"`
	HTTP             models.HTTP               `json:"http"`
	Auth             models.Auth               `json:"Auth"`
	DsAliCloudConfig models.DsAliCloudConfig   `json:"dsAliCloudConfig" `
	AWSCloudWatch    models.AWSCloudWatch      `json:"awsCloudwatch" `
	ClickHouseConfig models.DsClickHouseConfig `json:"clickhouseConfig"`
	Description      string                    `json:"description"`
	KubeConfig       string                    `json:"kubeConfig"`
	UpdateBy         string                    `json:"updateBy"`
	Enabled          *bool                     `json:"enabled" `
}

type RequestDatasourceUpdate struct {
	TenantId         string                    `json:"tenantId"`
	ID               string                    `json:"id"`
	Name             string                    `json:"name"`
	Labels           map[string]interface{}    `json:"labels" ` // 额外标签，会添加到事件Metric中，可用于区分数据来源；
	Type             string                    `json:"type"`
	HTTP             models.HTTP               `json:"http"`
	Auth             models.Auth               `json:"Auth"`
	DsAliCloudConfig models.DsAliCloudConfig   `json:"dsAliCloudConfig" `
	AWSCloudWatch    models.AWSCloudWatch      `json:"awsCloudwatch" `
	ClickHouseConfig models.DsClickHouseConfig `json:"clickhouseConfig"`
	Description      string                    `json:"description"`
	KubeConfig       string                    `json:"kubeConfig"`
	UpdateBy         string                    `json:"updateBy"`
	Enabled          *bool                     `json:"enabled" `
}

type RequestDatasourceQuery struct {
	TenantId string `json:"tenantId" form:"tenantId"`
	ID       string `json:"id" form:"id"`
	Type     string `json:"type" form:"type"`
	Query    string `json:"query" form:"query"`
}

type RequestQueryMetricsValue struct {
	DatasourceIds string `form:"datasourceIds"`
	Query         string `form:"query"`
	StartTime     int64  `form:"startTime"` // Unix 时间戳（秒），可选
	EndTime       int64  `form:"endTime"`   // Unix 时间戳（秒），可选
	Step          int64  `form:"step"`      // 步长（秒），可选
}

func (r RequestQueryMetricsValue) Validate() error {
	if r.Query == "" {
		return fmt.Errorf(ErrorQueryIsEmpty)
	}
	return nil
}

// GetStartTime 获取开始时间，如果未传则默认为过去 5 分钟
func (r RequestQueryMetricsValue) GetStartTime() time.Time {
	if r.StartTime == 0 {
		return time.Now().Add(-5 * time.Minute)
	}
	return time.Unix(r.StartTime, 0)
}

// GetEndTime 获取结束时间，如果未传则默认为当前时间
func (r RequestQueryMetricsValue) GetEndTime() time.Time {
	if r.EndTime == 0 {
		return time.Now()
	}
	return time.Unix(r.EndTime, 0)
}

// GetStep 获取步长，如果未传则默认为 10 秒
func (r RequestQueryMetricsValue) GetStep() time.Duration {
	if r.Step == 0 {
		return 10 * time.Second
	}
	return time.Duration(r.Step) * time.Second
}

type RequestSearchLogsContent struct {
	Type         string `json:"type"`
	DatasourceId string `json:"datasourceId"`
	Index        string `json:"index"`
	Query        string `json:"query"`
}

func (requestSearchLogsContent RequestSearchLogsContent) GetElasticSearchIndexName() string {
	if strings.Contains(requestSearchLogsContent.Index, "YYYY") && strings.Contains(requestSearchLogsContent.Index, "MM") && strings.Contains(requestSearchLogsContent.Index, "dd") {
		indexName := requestSearchLogsContent.Index
		indexName = strings.ReplaceAll(indexName, "YYYY", time.Now().Format("2006"))
		indexName = strings.ReplaceAll(indexName, "MM", time.Now().Format("01"))
		indexName = strings.ReplaceAll(indexName, "dd", time.Now().Format("02"))
		return indexName
	}

	return requestSearchLogsContent.Index
}
