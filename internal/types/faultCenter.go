package types

import "watchAlert/internal/models"

// RequestFaultCenterCreate 请求创建故障中心
type RequestFaultCenterCreate struct {
	TenantId              string                 `json:"tenantId"`
	Name                  string                 `json:"name"`
	Description           string                 `json:"description"`
	NoticeIds             []string               `json:"noticeIds" gorm:"column:noticeIds;serializer:json"`
	NoticeRoutes          []models.NoticeRoute   `json:"noticeRoutes" gorm:"noticeRoutes;serializer:json"`
	RepeatNoticeInterval  int64                  `json:"repeatNoticeInterval"`
	RecoverNotify         *bool                  `json:"recoverNotify"`
	AggregationType       string                 `json:"aggregationType"`
	CreateAt              int64                  `json:"createAt"`
	RecoverWaitTime       int64                  `json:"recoverWaitTime"` // 告警恢复等待时间，单位（秒）
	CurrentPreAlertNumber int64                  `json:"currentPreAlertNumber" gorm:"-"`
	CurrentAlertNumber    int64                  `json:"currentAlertNumber" gorm:"-"`
	CurrentMuteNumber     int64                  `json:"currentMuteNumber" gorm:"-"`
	CurrentRecoverNumber  int64                  `json:"currentRecoverNumber" gorm:"-"`
	IsUpgradeEnabled      *bool                  `json:"isUpgradeEnabled" gorm:"column:isUpgradeEnabled"`
	UpgradableSeverity    []string               `json:"upgradableSeverity" gorm:"column:upgradableSeverity;serializer:json"`
	UpgradeStrategy       models.UpgradeStrategy `json:"upgradeStrategy" gorm:"column:upgradeStrategy;serializer:json"`
}

// RequestFaultCenterUpdate 请求更新故障中心
type RequestFaultCenterUpdate struct {
	TenantId              string                 `json:"tenantId"`
	ID                    string                 `json:"id"`
	Name                  string                 `json:"name"`
	Description           string                 `json:"description"`
	NoticeIds             []string               `json:"noticeIds" gorm:"column:noticeIds;serializer:json"`
	NoticeRoutes          []models.NoticeRoute   `json:"noticeRoutes" gorm:"noticeRoutes;serializer:json"`
	RepeatNoticeInterval  int64                  `json:"repeatNoticeInterval"`
	RecoverNotify         *bool                  `json:"recoverNotify"`
	AggregationType       string                 `json:"aggregationType"`
	CreateAt              int64                  `json:"createAt"`
	RecoverWaitTime       int64                  `json:"recoverWaitTime"` // 告警恢复等待时间，单位（秒）
	CurrentPreAlertNumber int64                  `json:"currentPreAlertNumber" gorm:"-"`
	CurrentAlertNumber    int64                  `json:"currentAlertNumber" gorm:"-"`
	CurrentMuteNumber     int64                  `json:"currentMuteNumber" gorm:"-"`
	CurrentRecoverNumber  int64                  `json:"currentRecoverNumber" gorm:"-"`
	IsUpgradeEnabled      *bool                  `json:"isUpgradeEnabled" gorm:"column:isUpgradeEnabled"`
	UpgradableSeverity    []string               `json:"upgradableSeverity" gorm:"column:upgradableSeverity;serializer:json"`
	UpgradeStrategy       models.UpgradeStrategy `json:"upgradeStrategy" gorm:"column:upgradeStrategy;serializer:json"`
}

// RequestFaultCenterQuery 请求查询故障中心
type RequestFaultCenterQuery struct {
	TenantId string `form:"tenantId"`
	ID       string `form:"id"`
	Name     string `form:"name"`
	Query    string `from:"query"`
}

// RequestFaultCenterReset 请求重新配置故障中心
type RequestFaultCenterReset struct {
	TenantId        string `json:"tenantId"`
	ID              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	AggregationType string `json:"aggregationType"`
}

type RequestFaultCenterSLO struct {
	MTTA []float64 `json:"mtta"`
	MTTR []float64 `json:"mttr"`
}
