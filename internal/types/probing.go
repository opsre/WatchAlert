package types

import "watchAlert/internal/models"

// RequestProbingRuleCreate 请求创建拨测规则
type RequestProbingRuleCreate struct {
	TenantId              string                       `json:"tenantId"`
	RuleName              string                       `json:"ruleName"`
	RuleType              string                       `json:"ruleType"`
	RepeatNoticeInterval  int64                        `json:"repeatNoticeInterval"`
	ProbingEndpointConfig models.ProbingEndpointConfig `json:"probingEndpointConfig" `
	ProbingEndpointValues models.ProbingEndpointValues `json:"probingEndpointValues" `
	NoticeId              string                       `json:"noticeId"`
	Annotations           string                       `json:"annotations"`
	RecoverNotify         *bool                        `json:"recoverNotify"`
	UpdateAt              int64                        `json:"updateAt"`
	UpdateBy              string                       `json:"updateBy"`
	Enabled               *bool                        `json:"enabled" `
}

func (requestProbingRuleCreate *RequestProbingRuleCreate) GetEnabled() *bool {
	if requestProbingRuleCreate.Enabled == nil {
		isOk := false
		return &isOk
	}
	return requestProbingRuleCreate.Enabled
}

// RequestProbingRuleUpdate 请求更新拨测规则
type RequestProbingRuleUpdate struct {
	TenantId              string                       `json:"tenantId"`
	RuleId                string                       `json:"ruleId" gorm:"ruleId"`
	RuleName              string                       `json:"ruleName"`
	RuleType              string                       `json:"ruleType"`
	RepeatNoticeInterval  int64                        `json:"repeatNoticeInterval"`
	ProbingEndpointConfig models.ProbingEndpointConfig `json:"probingEndpointConfig" `
	ProbingEndpointValues models.ProbingEndpointValues `json:"probingEndpointValues" `
	NoticeId              string                       `json:"noticeId"`
	Annotations           string                       `json:"annotations"`
	RecoverNotify         *bool                        `json:"recoverNotify"`
	UpdateAt              int64                        `json:"updateAt"`
	UpdateBy              string                       `json:"updateBy"`
	Enabled               *bool                        `json:"enabled" `
}

func (requestProbingRuleUpdate *RequestProbingRuleUpdate) GetEnabled() *bool {
	if requestProbingRuleUpdate.Enabled == nil {
		isOk := false
		return &isOk
	}
	return requestProbingRuleUpdate.Enabled
}

// RequestProbingRuleQuery 请求查询拨测规则
type RequestProbingRuleQuery struct {
	TenantId string `json:"tenantId" form:"tenantId"`
	RuleId   string `json:"ruleId" form:"ruleId"`
	RuleType string `json:"ruleType" form:"ruleType"`
	Enabled  *bool  `json:"enabled" form:"enabled"`
	Query    string `json:"query" form:"query"`
}

// RequestProbingOnce 一次性拨测
type RequestProbingOnce struct {
	RuleType              string                       `json:"ruleType"`
	ProbingEndpointConfig models.ProbingEndpointConfig `json:"probingEndpointConfig"`
}

// RequestProbingHistoryRecord 获取历史拨测记录
type RequestProbingHistoryRecord struct {
	RuleId    string `json:"ruleId" form:"ruleId"`
	DateRange int64  `json:"dateRange" form:"dateRange"`
}

// RequestProbeChangeState 修改拨测规则状态
type RequestProbeChangeState struct {
	TenantId string `json:"tenantId"`
	RuleId   string `json:"ruleId"`
	Enabled  *bool  `json:"enabled"`
}

func (r *RequestProbeChangeState) GetEnabled() *bool {
	if r.Enabled == nil {
		isOk := false
		return &isOk
	}
	return r.Enabled
}
