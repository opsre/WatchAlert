package models

import (
	"fmt"
	"time"
)

// AlertStatus 定义状态类型
type AlertStatus string

// 所有可能的状态
const (
	StatePreAlert        AlertStatus = "pre_alert"        // 预告警
	StateAlerting        AlertStatus = "alerting"         // 告警中
	StatePendingRecovery AlertStatus = "pending_recovery" // 待恢复
	StateRecovered       AlertStatus = "recovered"        // 已恢复
	StateSilenced        AlertStatus = "silenced"         // 静默中
)

type AlertCurEvent struct {
	TenantId               string                 `json:"tenantId"`
	RuleId                 string                 `json:"rule_id"`
	RuleName               string                 `json:"rule_name"`
	DatasourceType         string                 `json:"datasource_type"`
	DatasourceId           string                 `json:"datasource_id" gorm:"datasource_id"`
	Fingerprint            string                 `json:"fingerprint"`
	Severity               string                 `json:"severity"`
	Labels                 map[string]interface{} `json:"labels" gorm:"labels;serializer:json"`
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
	UpgradeState           UpgradeState           `json:"upgradeState" gorm:"-"`
	Status                 AlertStatus            `json:"status" gorm:"-"` // 事件状态
}

type UpgradeState struct {
	IsConfirm       bool   `json:"isConfirm"`       // 是否已认领
	ConfirmOkTime   int64  `json:"confirmOkTime"`   // 点击认领时间
	ConfirmSendTime int64  `json:"confirmSendTime"` // 认领超时通知时间
	WhoAreConfirm   string `json:"whoAreConfirm"`

	IsHandle       bool   `json:"isHandle"`       // 是否已处理
	HandleOkTime   int64  `json:"HandleOkTime"`   // 点击处理时间
	HandleSendTime int64  `json:"handleSendTime"` // 处理超时通知时间
	WhoAreHandle   string `json:"whoAreHandle"`
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
	Status         string `json:"status" form:"status"`
	SortOrder      string `json:"sortOrder" form:"sortOrder"`
	Page
}

const (
	SortOrderASC  string = "ascend"
	SortOrderDesc string = "descend"
)

type ProcessAlertEvent struct {
	TenantId      string   `json:"tenantId"`
	State         int64    `json:"state"`
	FaultCenterId string   `json:"faultCenterId"`
	Fingerprints  []string `json:"fingerprints"`
	Time          int64    `json:"time"`
	Username      string   `json:"username"`
}

type CurEventResponse struct {
	List []AlertCurEvent `json:"list"`
	Page
}

func (alert *AlertCurEvent) TransitionStatus(newStatus AlertStatus) error {
	// 相同状态不需要转换
	if alert.Status == newStatus {
		return nil
	}

	// 检查当前状态是否允许转换到新状态
	if err := alert.validateTransition(newStatus); err != nil {
		return err
	}

	// 执行状态转换时的附加操作
	if err := alert.handleStateTransition(newStatus); err != nil {
		return err
	}

	// 更新状态
	alert.Status = newStatus

	return nil
}

// 验证状态转换是否有效
func (alert *AlertCurEvent) validateTransition(newState AlertStatus) error {
	current := alert.Status

	// 定义允许的状态转换规则
	allowedTransitions := map[AlertStatus][]AlertStatus{
		StatePreAlert:        {StateAlerting, StateSilenced},
		StateAlerting:        {StatePendingRecovery, StateSilenced},
		StatePendingRecovery: {StateAlerting, StateRecovered},
		StateRecovered:       {StatePreAlert},
		StateSilenced:        {StatePreAlert, StateAlerting, StatePendingRecovery, StateRecovered},
	}

	// 检查转换是否允许
	allowed := false
	for _, allowedState := range allowedTransitions[current] {
		if allowedState == newState {
			allowed = true
			break
		}
	}

	if !allowed {
		return StateTransitionError{
			FromState: current,
			ToState:   newState,
			Reason:    "不允许的状态转换",
		}
	}

	return nil
}

// handleStateTransition 处理状态转换时的附加操作
func (alert *AlertCurEvent) handleStateTransition(newState AlertStatus) error {
	now := time.Now().Unix()

	switch newState {
	case StatePreAlert:
		alert.FirstTriggerTime = now
		alert.LastEvalTime = now
	case StateAlerting:
	case StateRecovered:
		if alert.IsRecovered == true && alert.Status == StateRecovered {
			return nil
		}

		alert.LastSendTime = 0
		alert.RecoverTime = now
		alert.IsRecovered = true
	case StateSilenced:
	}

	return nil
}

// StateTransitionError 状态转换错误
type StateTransitionError struct {
	FromState AlertStatus
	ToState   AlertStatus
	Reason    string
}

func (e StateTransitionError) Error() string {
	return fmt.Sprintf("invalid transition from %s to %s: %s", e.FromState, e.ToState, e.Reason)
}

// IsArriveForDuration 比对持续时间
func (alert *AlertCurEvent) IsArriveForDuration() bool {
	return alert.LastEvalTime-alert.FirstTriggerTime > alert.ForDuration
}
