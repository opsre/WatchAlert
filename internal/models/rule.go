package models

type AlertRule struct {
	//gorm.Model
	TenantId             string            `json:"tenantId"`
	RuleId               string            `json:"ruleId" gorm:"ruleId"`
	RuleGroupId          string            `json:"ruleGroupId"`
	ExternalLabels       map[string]string `json:"externalLabels" gorm:"externalLabels;serializer:json"`
	DatasourceType       string            `json:"datasourceType"`
	DatasourceIdList     []string          `json:"datasourceId" gorm:"datasourceId;serializer:json"`
	RuleName             string            `json:"ruleName"`
	EvalInterval         int64             `json:"evalInterval"`
	EvalTimeType         string            `json:"evalTimeType"` // second, millisecond
	RepeatNoticeInterval int64             `json:"repeatNoticeInterval"`
	Description          string            `json:"description"`
	EffectiveTime        EffectiveTime     `json:"effectiveTime" gorm:"effectiveTime;serializer:json"`
	Severity             string            `json:"severity"`

	// Prometheus
	PrometheusConfig PrometheusConfig `json:"prometheusConfig" gorm:"prometheusConfig;serializer:json"`

	// 阿里云SLS
	AliCloudSLSConfig AliCloudSLSConfig `json:"alicloudSLSConfig" gorm:"alicloudSLSConfig;serializer:json"`

	// Loki
	LokiConfig LokiConfig `json:"lokiConfig" gorm:"lokiConfig;serializer:json"`

	VictoriaLogsConfig VictoriaLogsConfig `json:"victoriaLogsConfig" gorm:"victoriaConfig;serializer:json"`

	ClickHouseConfig ClickHouseConfig `json:"clickhouseConfig" gorm:"clickhouseConfig;serializer:json"`

	// Jaeger
	JaegerConfig JaegerConfig `json:"jaegerConfig" gorm:"JaegerConfig;serializer:json"`

	// AWS CloudWatch
	CloudWatchConfig CloudWatchConfig `json:"cloudwatchConfig" gorm:"cloudwatchConfig;serializer:json"`

	KubernetesConfig KubernetesConfig `json:"kubernetesConfig" gorm:"kubernetesConfig;serializer:json"`

	ElasticSearchConfig ElasticSearchConfig `json:"elasticSearchConfig" gorm:"elasticSearchConfig;serializer:json"`

	LogEvalCondition string `json:"logEvalCondition" gorm:"logEvalCondition;serializer:json"`

	FaultCenterId string `json:"faultCenterId"`
	Enabled       *bool  `json:"enabled" gorm:"enabled"`
}

type ElasticSearchConfig struct {
	Index           string            `json:"index"`
	Scope           int64             `json:"scope"`
	Filter          []EsQueryFilter   `json:"filter"`
	FilterCondition EsFilterCondition `json:"filterCondition"`
	EsQueryType     EsQueryType       `json:"queryType"`
	QueryWildcard   int64             `json:"queryWildcard"` // 0 精准匹配，1 模糊匹配
	RawJson         string            `json:"rawJson"`
}

type EsQueryType string

const (
	EsQueryTypeRawJson EsQueryType = "RawJson"
	EsQueryTypeField   EsQueryType = "Field"
)

type EsFilterCondition string

const (
	EsFilterConditionAnd EsFilterCondition = "And"
	EsFilterConditionOr  EsFilterCondition = "Or"
	EsFilterConditionNot EsFilterCondition = "Not"
)

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
	Project  string `json:"project"`
	Logstore string `json:"logstore"`
	LogQL    string `json:"logQL"`    // 查询语句
	LogScope int    `json:"logScope"` // 相对查询的日志范围（单位分钟）,1(min) 5(min)...
}

type LokiConfig struct {
	LogQL    string `json:"logQL"`
	LogScope int    `json:"logScope"`
}

type VictoriaLogsConfig struct {
	LogQL    string `json:"logQL"`
	LogScope int    `json:"logScope"`
	Limit    int    `json:"limit"`
}

type ClickHouseConfig struct {
	LogQL string `json:"logQL"`
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

// EvalCondition 评估表达式
type EvalCondition struct {
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

func (a *AlertRule) GetRuleType() string { return a.DatasourceType }

func (a *AlertRule) GetEnabled() *bool {
	if a.Enabled == nil {
		isOk := false
		return &isOk
	}
	return a.Enabled
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
