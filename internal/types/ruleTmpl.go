package types

import "watchAlert/internal/models"

type RequestRuleTemplateCreate struct {
	Type                 string                     `json:"type"`
	RuleGroupName        string                     `json:"ruleGroupName"`
	RuleName             string                     `json:"ruleName"`
	DatasourceType       string                     `json:"datasourceType"`
	EvalInterval         int64                      `json:"evalInterval"`
	ForDuration          int64                      `json:"forDuration"`
	RepeatNoticeInterval int64                      `json:"repeatNoticeInterval"`
	Description          string                     `json:"description"`
	PrometheusConfig     models.PrometheusConfig    `json:"prometheusConfig"`
	AliCloudSLSConfig    models.AliCloudSLSConfig   `json:"alicloudSLSConfig"`
	LokiConfig           models.LokiConfig          `json:"lokiConfig"`
	JaegerConfig         models.JaegerConfig        `json:"jaegerConfig"`
	KubernetesConfig     models.KubernetesConfig    `json:"kubernetesConfig"`
	ElasticSearchConfig  models.ElasticSearchConfig `json:"elasticSearchConfig"`
	VictoriaLogsConfig   models.VictoriaLogsConfig  `json:"victoriaLogsConfig"`
	ClickHouseConfig     models.ClickHouseConfig    `json:"clickhouseConfig"`
}

type RequestRuleTemplateUpdate struct {
	Type                 string                     `json:"type"`
	RuleGroupName        string                     `json:"ruleGroupName"`
	RuleName             string                     `json:"ruleName"`
	DatasourceType       string                     `json:"datasourceType"`
	EvalInterval         int64                      `json:"evalInterval"`
	ForDuration          int64                      `json:"forDuration"`
	RepeatNoticeInterval int64                      `json:"repeatNoticeInterval"`
	Description          string                     `json:"description"`
	PrometheusConfig     models.PrometheusConfig    `json:"prometheusConfig"`
	AliCloudSLSConfig    models.AliCloudSLSConfig   `json:"alicloudSLSConfig"`
	LokiConfig           models.LokiConfig          `json:"lokiConfig"`
	JaegerConfig         models.JaegerConfig        `json:"jaegerConfig"`
	KubernetesConfig     models.KubernetesConfig    `json:"kubernetesConfig"`
	ElasticSearchConfig  models.ElasticSearchConfig `json:"elasticSearchConfig"`
	VictoriaLogsConfig   models.VictoriaLogsConfig  `json:"victoriaLogsConfig"`
	ClickHouseConfig     models.ClickHouseConfig    `json:"clickhouseConfig"`
}

type RequestRuleTemplateQuery struct {
	Type           string `json:"type" form:"type"`
	RuleGroupName  string `json:"ruleGroupName" form:"ruleGroupName"`
	RuleName       string `json:"ruleName" form:"ruleName"`
	DatasourceType string `json:"datasourceType" form:"datasourceType"`
	Severity       int64  `json:"severity" form:"severity"`
	Annotations    string `json:"annotations" form:"annotations"`
	Query          string `json:"query" form:"query"`
	models.Page
}

type ResponseRuleTemplateList struct {
	List []models.RuleTemplate `json:"list"`
	models.Page
}
