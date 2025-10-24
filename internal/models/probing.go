package models

import "fmt"

type ProbingRule struct {
	TenantId              string                `json:"tenantId"`
	RuleName              string                `json:"ruleName"`
	RuleId                string                `json:"ruleId" gorm:"ruleId"`
	RuleType              string                `json:"ruleType"`
	RepeatNoticeInterval  int64                 `json:"repeatNoticeInterval"`
	ProbingEndpointConfig ProbingEndpointConfig `json:"probingEndpointConfig" gorm:"probingEndpointConfig;serializer:json"`
	ProbingEndpointValues ProbingEndpointValues `json:"probingEndpointValues" gorm:"-"`
	NoticeId              string                `json:"noticeId"`
	Annotations           string                `json:"annotations"`
	RecoverNotify         *bool                 `json:"recoverNotify"`
	UpdateAt              int64                 `json:"updateAt"`
	UpdateBy              string                `json:"updateBy"`
	Enabled               *bool                 `json:"enabled" gorm:"enabled"`
}

func (n *ProbingRule) TableName() string {
	return "w8t_probing_rule"
}

func (n *ProbingRule) GetRecoverNotify() *bool {
	if n.RecoverNotify == nil {
		isOk := false
		return &isOk
	}
	return n.RecoverNotify
}

func (n *ProbingRule) GetEnabled() *bool {
	if n.Enabled == nil {
		isOk := false
		return &isOk
	}
	return n.Enabled
}

type ProbingEndpointValues struct {
	PHTTP Phttp `json:"pHttp"`
	PICMP Picmp `json:"pIcmp"`
	PTCP  Ptcp  `json:"pTcp"`
	PSSL  Pssl  `json:"pSsl"`
}

type Picmp struct {
	// 丢包率的百分比
	PacketLoss string `json:"packetLoss"`
	// 最短的 RTT 时间, ms
	MinRtt string `json:"minRtt"`
	// 最长的 RTT 时间, ms
	MaxRtt string `json:"maxRtt"`
	// 平均 RTT 时间, ms
	AvgRtt string `json:"avgRtt"`
}

type Phttp struct {
	// 状态码
	StatusCode string `json:"statusCode" json:"status_code,omitempty"`
	// 响应时间, ms
	Latency string `json:"latency" json:"latency,omitempty"`
}

type Ptcp struct {
	IsSuccessful string `json:"isSuccessful"`
	ErrorMessage string `json:"errorMessage"`
}

type Pssl struct {
	ExpireTime    string `json:"expireTime"`
	ResponseTime  string `json:"responseTime"`
	StartTime     string `json:"startTime"`
	TimeRemaining string `json:"timeRemaining"`
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
	// 失败次数
	Failure int `json:"failure"`
	// 运算
	Operator string `json:"operator"`
	// 字段
	Field string `json:"field"`
	// 预期值
	ExpectedValue float64 `json:"expectedValue"`
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

// ------------------------ Event ------------------------

type ProbingEvent struct {
	TenantId               string                 `json:"tenantId"`
	RuleId                 string                 `json:"ruleId" gorm:"ruleId"`
	RuleName               string                 `json:"ruleName"`
	RuleType               string                 `json:"ruleType"`
	Fingerprint            string                 `json:"fingerprint"`
	Labels                 map[string]interface{} `json:"labels" gorm:"labels;serializer:json"`
	ProbingEndpointConfig  ProbingEndpointConfig  `json:"probingEndpointConfig" gorm:"probingEndpointConfig;serializer:json"`
	NoticeId               string                 `json:"noticeId"`
	IsRecovered            bool                   `json:"isRecovered" gorm:"-"`
	RecoverNotify          *bool                  `json:"recoverNotify"`
	FirstTriggerTime       int64                  `json:"first_trigger_time"` // 第一次触发时间
	FirstTriggerTimeFormat string                 `json:"first_trigger_time_format" gorm:"-"`
	RepeatNoticeInterval   int64                  `json:"repeat_notice_interval"`  // 重复通知间隔时间
	LastEvalTime           int64                  `json:"last_eval_time" gorm:"-"` // 上一次评估时间
	LastSendTime           int64                  `json:"last_send_time" gorm:"-"` // 上一次发送时间
	RecoverTime            int64                  `json:"recover_time" gorm:"-"`   // 恢复时间
	DutyUser               string                 `json:"duty_user" gorm:"-"`
	RecoverTimeFormat      string                 `json:"recover_time_format" gorm:"-"`
	Annotations            string                 `json:"annotations" gorm:"-"`
}

func (n *ProbingEvent) GetRecoverNotify() *bool {
	if n.RecoverNotify == nil {
		isOk := false
		return &isOk
	}
	return n.RecoverNotify
}

type ProbingEventCacheKey string

func BuildProbingEventCacheKey(tenantId, ruleId string) ProbingEventCacheKey {
	return ProbingEventCacheKey(fmt.Sprintf("w8t:%s:probing:%s.event", tenantId, ruleId))
}

type ProbingValueCacheKey string

func BuildProbingValueCacheKey(tenantId, ruleId string) ProbingValueCacheKey {
	return ProbingValueCacheKey(fmt.Sprintf("w8t:%s:probing:%s.value", tenantId, ruleId))
}

type ProbingHistory struct {
	Timestamp int64          `json:"timestamp"`
	RuleId    string         `json:"ruleId"`
	Value     map[string]any `json:"value" gorm:"value;serializer:json"`
}

func (p *ProbingHistory) TableName() string {
	return "w8t_probing_history"
}
