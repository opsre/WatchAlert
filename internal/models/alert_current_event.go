package models

import (
	"fmt"
)

type AlertCurEvent struct {
	TenantId               string                 `json:"tenantId"`
	RuleId                 string                 `json:"rule_id"`
	RuleName               string                 `json:"rule_name"`
	DatasourceType         string                 `json:"datasource_type"`
	DatasourceId           string                 `json:"datasource_id" gorm:"datasource_id"`
	Fingerprint            string                 `json:"fingerprint"`
	Severity               string                 `json:"severity"`
	Metric                 map[string]interface{} `json:"metric" gorm:"metric;serializer:json"`
	Labels                 map[string]string      `json:"labels" gorm:"labels;serializer:json"`
	SearchQL               string                 `json:"searchQL" gorm:"-"`
	EvalInterval           int64                  `json:"eval_interval"`
	ForDuration            int64                  `json:"for_duration"`
	Annotations            string                 `json:"annotations" gorm:"-"`
	IsRecovered            bool                   `json:"is_recovered" gorm:"-"`
	FirstTriggerTime       int64                  `json:"first_trigger_time"` // 第一次触发时间
	FirstTriggerTimeFormat string                 `json:"first_trigger_time_format" gorm:"-"`
	RepeatNoticeInterval   int64                  `json:"repeat_notice_interval"`  // 重复通知间隔时间
	LastEvalTime           int64                  `json:"last_eval_time" gorm:"-"` // 上一次评估时间
	LastSendTime           int64                  `json:"last_send_time" gorm:"-"` // 上一次发送时间
	RecoverTime            int64                  `json:"recover_time" gorm:"-"`   // 恢复时间
	RecoverTimeFormat      string                 `json:"recover_time_format" gorm:"-"`
	DutyUser               string                 `json:"duty_user" gorm:"-"`
	DutyUserPhoneNumber    []string               `json:"duty_user_phone_number" gorm:"-"`
	EffectiveTime          EffectiveTime          `json:"effectiveTime" gorm:"effectiveTime;serializer:json"`
	FaultCenterId          string                 `json:"faultCenterId"`
	FaultCenter            FaultCenter            `json:"faultCenter" gorm:"-"`
	Status                 int64                  `json:"status" gorm:"-"` // 事件状态，预告警：0，告警中：1，静默中：2，待恢复：3，已恢复：4
}

type AlertCurEventQuery struct {
	TenantId       string `json:"tenantId" form:"tenantId"`
	RuleId         string `json:"ruleId" form:"ruleId"`
	RuleName       string `json:"ruleName" form:"ruleName"`
	DatasourceType string `json:"datasourceType" form:"datasourceType"`
	DatasourceId   string `json:"datasourceId" form:"datasourceId"`
	Fingerprint    string `json:"fingerprint" form:"fingerprint"`
	Query          string `json:"query" form:"query"`
	Scope          int64  `json:"scope" form:"scope"`
	Severity       string `json:"severity" form:"severity"`
	FaultCenterId  string `json:"faultCenterId" form:"faultCenterId"`
	Page
}

type CurEventResponse struct {
	List []AlertCurEvent `json:"list"`
	Page
}

func (ace *AlertCurEvent) GetCacheEventsKey() string {
	return fmt.Sprintf("w8t:%s:%s:%s.events", ace.TenantId, FaultCenterPrefix, ace.FaultCenterId)
}

// IsArriveForDuration 比对持续时间
func (ace *AlertCurEvent) IsArriveForDuration() bool {
	return ace.LastEvalTime-ace.FirstTriggerTime > ace.ForDuration
}

// DetermineEventStatus 判断事件状态
func (ace *AlertCurEvent) DetermineEventStatus() int64 {
	if !ace.IsArriveForDuration() {
		return 0 // 未达到持续时间
	}

	if ace.IsRecovered {
		return 3 // 事件待恢复
	}

	return 1 // 事件处于触发状态
}
