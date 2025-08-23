package types

import (
	"strings"
	"time"
	"watchAlert/internal/models"
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
