package models

import (
	"fmt"
	"slices"
)

// 常量定义
const (
	FaultCenterPrefix = "faultCenter"
	ConfirmStatus     = 1
	HandleStatus      = 2
)

type FaultCenter struct {
	TenantId              string            `json:"tenantId"`
	ID                    string            `json:"id"`
	Name                  string            `json:"name"`
	Description           string            `json:"description"`
	NoticeIds             []string          `json:"noticeIds" gorm:"column:noticeIds;serializer:json"`
	NoticeRoutes          []NoticeRoute     `json:"noticeRoutes" gorm:"noticeRoutes;serializer:json"`
	RepeatNoticeInterval  int64             `json:"repeatNoticeInterval"`
	RecoverNotify         *bool             `json:"recoverNotify"`
	AggregationType       string            `json:"aggregationType"`
	CreateAt              int64             `json:"createAt"`
	RecoverWaitTime       int64             `json:"recoverWaitTime"`
	CurrentPreAlertNumber int64             `json:"currentPreAlertNumber" gorm:"-"`
	CurrentAlertNumber    int64             `json:"currentAlertNumber" gorm:"-"`
	CurrentMuteNumber     int64             `json:"currentMuteNumber" gorm:"-"`
	CurrentRecoverNumber  int64             `json:"currentRecoverNumber" gorm:"-"`
	IsUpgradeEnabled      *bool             `json:"isUpgradeEnabled" gorm:"column:isUpgradeEnabled"`
	UpgradableSeverity    []string          `json:"upgradableSeverity" gorm:"column:upgradableSeverity;serializer:json"`
	UpgradeStrategy       []UpgradeStrategy `json:"upgradeStrategy" gorm:"column:upgradeStrategy;serializer:json"`
}

type UpgradeStrategy struct {
	Enabled        *bool  `json:"enabled"`
	StrategyType   int64  `json:"strategyType"`   // 策略类型，1 认领超时 / 2 处理超时
	Timeout        int64  `json:"timeout"`        // 超时时间
	RepeatInterval int64  `json:"repeatInterval"` // 重复通知间隔时间
	NoticeId       string `json:"noticeId"`       // 通知对象ID
}

func (u *UpgradeStrategy) GetEnabled() bool {
	if u.Enabled == nil {
		return false
	}
	return *u.Enabled
}

// GetSeverityAssessmentResult 获取等级评估结果，不满足条件时不进行升级
func (f *FaultCenter) GetSeverityAssessmentResult(severity string) bool {
	return slices.Contains(f.UpgradableSeverity, severity)
}

func (f *FaultCenter) GetNoticeInterval(strategyType int64) int64 {
	if s := f.GetStrategy(strategyType); s != nil {
		return s.RepeatInterval
	}
	return 0
}

func (f *FaultCenter) GetTimeout(strategyType int64) int64 {
	if s := f.GetStrategy(strategyType); s != nil {
		return s.Timeout
	}
	return 0
}

// GetUpgradeNoticeId 获取通知对象Id
func (f *FaultCenter) GetUpgradeNoticeId(strategyType int64) string {
	if s := f.GetStrategy(strategyType); s != nil {
		return s.NoticeId
	}
	return ""
}

func (f *FaultCenter) GetStrategy(strategyType int64) *UpgradeStrategy {
	for _, s := range f.UpgradeStrategy {
		if s.StrategyType != strategyType {
			continue
		}

		return &s
	}
	return nil
}

type NoticeRoute struct {
	Key       string   `json:"key"`
	Value     string   `json:"value"`
	NoticeIds []string `json:"noticeIds" gorm:"column:noticeIds;serializer:json"`
}

func (f *FaultCenter) TableName() string {
	return "w8t_fault_center"
}

func (f *FaultCenter) GetIsUpgradeEnabled() bool {
	if f.IsUpgradeEnabled == nil {
		return false
	}
	return *f.IsUpgradeEnabled
}

func (f *FaultCenter) GetRecoverNotify() bool {
	if f.RecoverNotify == nil {
		return false
	}
	return *f.RecoverNotify
}

func (f *FaultCenter) GetAlarmAggregationType() string {
	return f.AggregationType
}

type FaultCenterQuery struct {
	TenantId string `form:"tenantId"`
	ID       string `form:"id"`
	Name     string `form:"name"`
	Query    string `from:"query"`
}

type AlertEventCacheKey string

func BuildAlertEventCacheKey(tenantId, faultCenterId string) AlertEventCacheKey {
	return AlertEventCacheKey(fmt.Sprintf("w8t:%s:%s:%s.events", tenantId, FaultCenterPrefix, faultCenterId))
}

type AlertMuteCacheKey string

func BuildAlertMuteCacheKey(tenantId, faultCenterId string) AlertMuteCacheKey {
	return AlertMuteCacheKey(fmt.Sprintf("w8t:%s:%s:%s.mutes", tenantId, FaultCenterPrefix, faultCenterId))
}

type FaultCenterInfoCacheKey string

func BuildFaultCenterInfoCacheKey(tenantId, faultCenterId string) FaultCenterInfoCacheKey {
	return FaultCenterInfoCacheKey(fmt.Sprintf("w8t:%s:%s:%s.info", tenantId, FaultCenterPrefix, faultCenterId))
}
