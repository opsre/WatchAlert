package types

import "watchAlert/internal/models"

type RequestRuleCreate struct {
	TenantId             string                     `json:"tenantId"`
	RuleGroupId          string                     `json:"ruleGroupId"`
	ExternalLabels       map[string]string          `json:"externalLabels"`
	DatasourceType       string                     `json:"datasourceType"`
	DatasourceIdList     []string                   `json:"datasourceId"`
	RuleName             string                     `json:"ruleName"`
	EvalInterval         int64                      `json:"evalInterval"`
	EvalTimeType         string                     `json:"evalTimeType"` // second, millisecond
	RepeatNoticeInterval int64                      `json:"repeatNoticeInterval"`
	Description          string                     `json:"description"`
	EffectiveTime        models.EffectiveTime       `json:"effectiveTime"`
	Severity             string                     `json:"severity"`
	PrometheusConfig     models.PrometheusConfig    `json:"prometheusConfig"`
	AliCloudSLSConfig    models.AliCloudSLSConfig   `json:"alicloudSLSConfig"`
	LokiConfig           models.LokiConfig          `json:"lokiConfig"`
	VictoriaLogsConfig   models.VictoriaLogsConfig  `json:"victoriaLogsConfig"`
	ClickHouseConfig     models.ClickHouseConfig    `json:"clickhouseConfig"`
	JaegerConfig         models.JaegerConfig        `json:"jaegerConfig"`
	CloudWatchConfig     models.CloudWatchConfig    `json:"cloudwatchConfig"`
	KubernetesConfig     models.KubernetesConfig    `json:"kubernetesConfig"`
	ElasticSearchConfig  models.ElasticSearchConfig `json:"elasticSearchConfig"`
	LogEvalCondition     string                     `json:"logEvalCondition"`
	FaultCenterId        string                     `json:"faultCenterId"`
	UpdateBy             string                     `json:"updateBy"`
	Enabled              *bool                      `json:"enabled"`
}

func (requestRuleCreate *RequestRuleCreate) GetEnabled() *bool {
	if requestRuleCreate.Enabled == nil {
		isOk := false
		return &isOk
	}
	return requestRuleCreate.Enabled
}

type RequestRuleUpdate struct {
	TenantId             string                     `json:"tenantId"`
	RuleId               string                     `json:"ruleId"`
	RuleGroupId          string                     `json:"ruleGroupId"`
	ExternalLabels       map[string]string          `json:"externalLabels"`
	DatasourceType       string                     `json:"datasourceType"`
	DatasourceIdList     []string                   `json:"datasourceId"`
	RuleName             string                     `json:"ruleName"`
	EvalInterval         int64                      `json:"evalInterval"`
	EvalTimeType         string                     `json:"evalTimeType"` // second, millisecond
	RepeatNoticeInterval int64                      `json:"repeatNoticeInterval"`
	Description          string                     `json:"description"`
	EffectiveTime        models.EffectiveTime       `json:"effectiveTime"`
	Severity             string                     `json:"severity"`
	PrometheusConfig     models.PrometheusConfig    `json:"prometheusConfig"`
	AliCloudSLSConfig    models.AliCloudSLSConfig   `json:"alicloudSLSConfig"`
	LokiConfig           models.LokiConfig          `json:"lokiConfig"`
	VictoriaLogsConfig   models.VictoriaLogsConfig  `json:"victoriaLogsConfig"`
	ClickHouseConfig     models.ClickHouseConfig    `json:"clickhouseConfig"`
	JaegerConfig         models.JaegerConfig        `json:"jaegerConfig"`
	CloudWatchConfig     models.CloudWatchConfig    `json:"cloudwatchConfig"`
	KubernetesConfig     models.KubernetesConfig    `json:"kubernetesConfig"`
	ElasticSearchConfig  models.ElasticSearchConfig `json:"elasticSearchConfig"`
	LogEvalCondition     string                     `json:"logEvalCondition"`
	FaultCenterId        string                     `json:"faultCenterId"`
	UpdateBy             string                     `json:"updateBy"`
	Enabled              *bool                      `json:"enabled"`
}

func (requestRuleUpdate *RequestRuleUpdate) GetEnabled() *bool {
	if requestRuleUpdate.Enabled == nil {
		isOk := false
		return &isOk
	}
	return requestRuleUpdate.Enabled
}

type RequestRuleQuery struct {
	TenantId         string   `json:"tenantId" form:"tenantId"`
	RuleId           string   `json:"ruleId" form:"ruleId"`
	RuleGroupId      string   `json:"ruleGroupId" form:"ruleGroupId"`
	DatasourceType   string   `json:"datasourceType" form:"datasourceType"`
	DatasourceIdList []string `json:"datasourceId" form:"datasourceId"`
	RuleName         string   `json:"ruleName" form:"ruleName"`
	Enabled          string   `json:"enabled" form:"enabled"`
	Query            string   `json:"query" form:"query"`
	Status           string   `json:"status" form:"status"` // 查询规则状态
	models.Page
}

type ResponseRuleList struct {
	List []models.AlertRule `json:"list"`
	models.Page
}

type RequestRuleChangeStatus struct {
	TenantId      string `json:"tenantId" form:"tenantId"`
	RuleId        string `json:"ruleId" form:"ruleId"`
	RuleGroupId   string `json:"ruleGroupId" form:"ruleGroupId"`
	FaultCenterId string `json:"faultCenterId" form:"faultCenterId"`
	Enabled       *bool  `json:"enabled" form:"enabled"`
}

func (r *RequestRuleChangeStatus) GetEnabled() *bool {
	if r.Enabled == nil {
		isOk := false
		return &isOk
	}
	return r.Enabled
}

const (
	WithPrometheusRuleImport int = 0
	WithWatchAlertJsonImport int = 1
)

type RequestRuleImport struct {
	TenantId         string   `json:"tenantId"`
	RuleGroupId      string   `json:"ruleGroupId"`
	DatasourceType   string   `json:"datasourceType"`
	DatasourceIdList []string `json:"datasourceIdList"`
	FaultCenterId    string   `json:"faultCenterId"`
	ImportType       int      `json:"importType"`
	Rules            string   `json:"rules"`
}

type PrometheusAlerts struct {
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	Alert       string            `yaml:"alert"`
	Expr        string            `yaml:"expr"`
	For         string            `yaml:"for"`
	Labels      map[string]string `yaml:"labels"`
	Annotations Annotations       `yaml:"annotations"`
}

type Annotations struct {
	Description string `yaml:"description"`
}

func (r Rule) GetEnable() *bool {
	var enable = false
	return &enable
}
