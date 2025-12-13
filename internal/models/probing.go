package models

type ProbeRule struct {
	TenantId              string                `json:"tenantId"`
	RuleName              string                `json:"ruleName"`
	RuleId                string                `json:"ruleId" gorm:"ruleId"`
	RuleType              string                `json:"ruleType"`
	ProbingEndpointConfig ProbingEndpointConfig `json:"probingEndpointConfig" gorm:"probingEndpointConfig;serializer:json"`
	DatasourceId          string                `json:"datasourceId"`
	UpdateAt              int64                 `json:"updateAt"`
	UpdateBy              string                `json:"updateBy"`
	Enabled               *bool                 `json:"enabled" gorm:"enabled"`
}

func (n *ProbeRule) TableName() string {
	return "w8t_probe_rule"
}

func (n *ProbeRule) GetEnabled() *bool {
	if n.Enabled == nil {
		isOk := false
		return &isOk
	}
	return n.Enabled
}

type ProbingEndpointConfig struct {
	// 端点
	Endpoint string `json:"endpoint"`
	// 评估策略
	Strategy endpointStrategy `json:"strategy"`
	HTTP     ehttp            `json:"http"`
	ICMP     eicmp            `json:"icmp"`
}

type endpointStrategy struct {
	// 超时时间
	Timeout int `json:"timeout"`
	// 执行频率
	EvalInterval int64 `json:"evalInterval"`
}

type ehttp struct {
	Method string            `json:"method"`
	Header map[string]string `json:"header"`
	Body   string            `json:"body"`
}

type eicmp struct {
	Interval int `json:"interval"`
	Count    int `json:"count"`
}
