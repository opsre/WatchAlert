package models

import (
	"sort"
	"watchAlert/pkg/tools"
)

type Duration int64

type LabelsMap map[string]string

type NoticeGroup []map[string]string

type AlertRule struct {
	//gorm.Model
	TenantId             string        `json:"tenantId"`
	RuleId               string        `json:"ruleId" gorm:"ruleId"`
	RuleGroupId          string        `json:"ruleGroupId"`
	DatasourceType       string        `json:"datasourceType"`
	DatasourceIdList     []string      `json:"datasourceId" gorm:"datasourceId;serializer:json"`
	RuleName             string        `json:"ruleName"`
	EvalInterval         int64         `json:"evalInterval"`
	RepeatNoticeInterval int64         `json:"repeatNoticeInterval"`
	Description          string        `json:"description"`
	Labels               LabelsMap     `json:"labels" gorm:"labels;serializer:json"`
	EffectiveTime        EffectiveTime `json:"effectiveTime" gorm:"effectiveTime;serializer:json"`
	Severity             string        `json:"severity"`

	// Prometheus
	PrometheusConfig PrometheusConfig `json:"prometheusConfig" gorm:"prometheusConfig;serializer:json"`

	// 阿里云SLS
	AliCloudSLSConfig AliCloudSLSConfig `json:"alicloudSLSConfig" gorm:"alicloudSLSConfig;serializer:json"`

	// Loki
	LokiConfig LokiConfig `json:"lokiConfig" gorm:"lokiConfig;serializer:json"`

	// Jaeger
	JaegerConfig JaegerConfig `json:"jaegerConfig" gorm:"JaegerConfig;serializer:json"`

	// AWS CloudWatch
	CloudWatchConfig CloudWatchConfig `json:"cloudwatchConfig" gorm:"cloudwatchConfig;serializer:json"`

	KubernetesConfig KubernetesConfig `json:"kubernetesConfig" gorm:"kubernetesConfig;serializer:json"`

	ElasticSearchConfig ElasticSearchConfig `json:"elasticSearchConfig" gorm:"elasticSearchConfig;serializer:json"`

	NetworkEndpointConfig ProbingEndpointConfig `json:"networkEndpointConfig" gorm:"networkEndpointConfig;serializer:json"`

	NoticeId         string      `json:"noticeId"`
	NoticeGroup      NoticeGroup `json:"noticeGroup" gorm:"noticeGroup;serializer:json"`
	RecoverNotify    *bool       `json:"recoverNotify"`
	AlarmAggregation *bool       `json:"alarmAggregation"`
	Enabled          *bool       `json:"enabled" gorm:"enabled"`
}

type ElasticSearchConfig struct {
	Index  string          `json:"index"`
	Scope  int64           `json:"scope"`
	Filter []EsQueryFilter `json:"filter"`
}

type EsQueryFilter struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type KubernetesConfig struct {
	Resource string   `json:"resource"`
	Reason   string   `json:"reason"`
	Value    int      `json:"value"`
	Filter   []string `json:"filter"`
	Scope    int      `json:"scope"`
}

type JaegerConfig struct {
	Service string `json:"service"`
	Scope   int    `json:"scope"`
	Tags    string `json:"tags"`
}

type PrometheusConfig struct {
	PromQL      string  `json:"promQL"`
	Annotations string  `json:"annotations"`
	ForDuration int64   `json:"forDuration"`
	Rules       []Rules `json:"rules"`
}

type Rules struct {
	Severity string `json:"severity"`
	Expr     string `json:"expr"`
}

type EffectiveTime struct {
	Week      []string `json:"week"`
	StartTime int      `json:"startTime"`
	EndTime   int      `json:"endTime"`
}

type AliCloudSLSConfig struct {
	Project       string        `json:"project"`
	Logstore      string        `json:"logstore"`
	LogQL         string        `json:"logQL"`    // 查询语句
	LogScope      int           `json:"logScope"` // 相对查询的日志范围（单位分钟）,1(min) 5(min)...
	EvalCondition EvalCondition `json:"evalCondition" gorm:"evalCondition;serializer:json"`
}

type LokiConfig struct {
	LogQL         string        `json:"logQL"`
	LogScope      int           `json:"logScope"`
	EvalCondition EvalCondition `json:"evalCondition" gorm:"evalCondition;serializer:json"`
}

type CloudWatchConfig struct {
	Namespace  string   `json:"namespace"`
	MetricName string   `json:"metricName"`
	Statistic  string   `json:"statistic"`
	Period     int      `json:"period"`
	Expr       string   `json:"expr"`
	Threshold  int      `json:"threshold"`
	Dimension  string   `json:"dimension"`
	Endpoints  []string `json:"endpoints" gorm:"endpoints;serializer:json"`
}

// EvalCondition 日志评估条件
type EvalCondition struct {
	Type string `json:"type"`
	// 运算
	Operator string `json:"operator"`
	// 查询值
	QueryValue float64 `json:"queryValue"`
	// 预期值
	ExpectedValue float64 `json:"value"`
}

type Fingerprint uint64

type AlertRuleQuery struct {
	TenantId         string   `json:"tenantId" form:"tenantId"`
	RuleId           string   `json:"ruleId" form:"ruleId"`
	RuleGroupId      string   `json:"ruleGroupId" form:"ruleGroupId"`
	DatasourceType   string   `json:"datasourceType" form:"datasourceType"`
	DatasourceIdList []string `json:"datasourceId" form:"datasourceId"`
	RuleName         string   `json:"ruleName" form:"ruleName"`
	Enabled          string   `json:"enabled" form:"enabled"`
	Query            string   `json:"query" form:"query"`
	Status           string   `json:"status" form:"status"` // 查询规则状态
	Page
}

type RuleResponse struct {
	List []AlertRule `json:"list"`
	Page
}

var (
	// cache the signature of an empty label set.
	emptyLabelSignature = tools.HashNew()
)

const SeparatorByte byte = 255

// Fingerprint returns a unique hash for the alert. It is equivalent to
// the fingerprint of the alert's label set.
func (a *AlertRule) Fingerprint() Fingerprint {

	// 没有配置标签，则用随机生成
	if len(a.Labels) == 0 {
		return Fingerprint(emptyLabelSignature)
	}

	// 定义map存储所有标签
	labelNames := make([]string, 0, len(a.Labels))
	for labelName := range a.Labels {
		labelNames = append(labelNames, labelName)
	}
	// 标签排序。用于根据标签做hash
	sort.Strings(labelNames)

	// 在随机生成的hash的基础上，新增标签hash
	sum := tools.HashNew()
	for _, labelName := range labelNames {
		sum = tools.HashAdd(sum, labelName)
		sum = tools.HashAddByte(sum, SeparatorByte)
		sum = tools.HashAdd(sum, a.Labels[labelName])
		sum = tools.HashAddByte(sum, SeparatorByte)
	}
	return Fingerprint(sum)

}

func (a *AlertRule) GetRuleType() string { return a.DatasourceType }

func (a *AlertRule) GetRecoverNotify() *bool {
	if a.RecoverNotify == nil {
		isOk := false
		return &isOk
	}
	return a.RecoverNotify
}

func (a *AlertRule) GetAlarmAggregation() *bool {
	if a.AlarmAggregation == nil {
		isOk := false
		return &isOk
	}
	return a.AlarmAggregation
}

func (a *AlertRule) GetEnabled() *bool {
	if a.Enabled == nil {
		isOk := false
		return &isOk
	}
	return a.RecoverNotify
}
