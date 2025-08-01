package models

type RuleTemplateGroup struct {
	Name        string `json:"name" gorm:"type:varchar(255);not null"`
	Number      int    `json:"number"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type RuleTemplate struct {
	Type                 string              `json:"type"`
	RuleGroupName        string              `json:"ruleGroupName"`
	RuleName             string              `json:"ruleName"  gorm:"type:varchar(255);not null"`
	DatasourceType       string              `json:"datasourceType"`
	EvalInterval         int64               `json:"evalInterval"`
	ForDuration          int64               `json:"forDuration"`
	RepeatNoticeInterval int64               `json:"repeatNoticeInterval"`
	Description          string              `json:"description"`
	PrometheusConfig     PrometheusConfig    `json:"prometheusConfig" gorm:"prometheusConfig;serializer:json"`
	AliCloudSLSConfig    AliCloudSLSConfig   `json:"alicloudSLSConfig" gorm:"alicloudSLSConfig;serializer:json"`
	LokiConfig           LokiConfig          `json:"lokiConfig" gorm:"lokiConfig;serializer:json"`
	JaegerConfig         JaegerConfig        `json:"jaegerConfig" gorm:"JaegerConfig;serializer:json"`
	KubernetesConfig     KubernetesConfig    `json:"kubernetesConfig" gorm:"kubernetesConfig;serializer:json"`
	ElasticSearchConfig  ElasticSearchConfig `json:"elasticSearchConfig" gorm:"elasticSearchConfig;serializer:json"`
	VictoriaLogsConfig   VictoriaLogsConfig  `json:"victoriaLogsConfig" gorm:"victoriaConfig;serializer:json"`
	ClickHouseConfig     ClickHouseConfig    `json:"clickhouseConfig" gorm:"clickhouseConfig;serializer:json"`
}
