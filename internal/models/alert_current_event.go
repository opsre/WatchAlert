package models

import (
	"fmt"
	"time"
	"watchAlert/pkg/tools"
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
	EventId                string                 `json:"eventId"`
	RuleGroupId            string                 `json:"rule_group_id"`
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
	ConfirmState           ConfirmState           `json:"confirmState" gorm:"-"`
	Status                 AlertStatus            `json:"status" gorm:"-"` // 事件状态
}

type ConfirmState struct {
	IsOk                   bool   `json:"isOk"`                   // 是否已认领
	ConfirmActionTime      int64  `json:"confirmActionTime"`      // 点击认领时间
	ConfirmTimeoutSendTime int64  `json:"confirmTimeoutSendTime"` // 认领超时通知时间
	ConfirmUsername        string `json:"confirmUsername"`
}

const (
	SortOrderASC  string = "ascend"
	SortOrderDesc string = "descend"
)

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

// GetLastSendTime 获取故障中心事件的最后发送时间
func (alert *AlertCurEvent) GetLastSendTime() int64 {
	return alert.LastSendTime
}

// GetLastEvalTime 获取故障中心事件的最后评估时间
func (alert *AlertCurEvent) GetLastEvalTime() int64 {
	return time.Now().Unix()
}

// GetFirstTime 获取故障中心事件的首次触发时间
func (alert *AlertCurEvent) GetFirstTime() int64 {
	if alert.FirstTriggerTime == 0 {
		return time.Now().Unix()
	}
	return alert.FirstTriggerTime
}

// GetLastConfirmState 获取最新告警升级认领状态
func (alert *AlertCurEvent) GetLastConfirmState() ConfirmState {
	return alert.ConfirmState
}

// GetEventStatus 获取事件状态
func (alert *AlertCurEvent) GetEventStatus() AlertStatus {
	if alert.Status == "" {
		return StatePreAlert
	}
	return alert.Status
}

// GetEventId 获取告警事件ID
func (alert *AlertCurEvent) GetEventId() string {
	if alert.EventId == "" {
		return tools.RandId()
	}
	return alert.EventId
}
