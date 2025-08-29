package models

type RuleTemplateGroup struct {
	Name        string `json:"name" gorm:"type:varchar(255);not null"`
	Number      int    `json:"number"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type RuleTemplateGroupQuery struct {
	Name        string `json:"name" form:"name"`
	Type        string `json:"type" form:"type"`
	Description string `json:"description" form:"description"`
	Query       string `json:"query" form:"query"`
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
	EffectiveTime        EffectiveTime       `json:"effectiveTime" gorm:"effectiveTime;serializer:json"`
	PrometheusConfig     PrometheusConfig    `json:"prometheusConfig" gorm:"prometheusConfig;serializer:json"`
	AliCloudSLSConfig    AliCloudSLSConfig   `json:"alicloudSLSConfig" gorm:"alicloudSLSConfig;serializer:json"`
	LokiConfig           LokiConfig          `json:"lokiConfig" gorm:"lokiConfig;serializer:json"`
	JaegerConfig         JaegerConfig        `json:"jaegerConfig" gorm:"JaegerConfig;serializer:json"`
	KubernetesConfig     KubernetesConfig    `json:"kubernetesConfig" gorm:"kubernetesConfig;serializer:json"`
	ElasticSearchConfig  ElasticSearchConfig `json:"elasticSearchConfig" gorm:"elasticSearchConfig;serializer:json"`
}

type RuleTemplateQuery struct {
	Type           string `json:"type" form:"type"`
	RuleGroupName  string `json:"ruleGroupName" form:"ruleGroupName"`
	RuleName       string `json:"ruleName" form:"ruleName"`
	DatasourceType string `json:"datasourceType" form:"datasourceType"`
	Severity       int64  `json:"severity" form:"severity"`
	Annotations    string `json:"annotations" form:"annotations"`
	Query          string `json:"query" form:"query"`
}
