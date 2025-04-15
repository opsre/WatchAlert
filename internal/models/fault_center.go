package models

import "fmt"

const FaultCenterPrefix = "faultCenter"

type FaultCenter struct {
	TenantId              string        `json:"tenantId"`
	ID                    string        `json:"id"`
	Name                  string        `json:"name"`
	Description           string        `json:"description"`
	NoticeIds             []string      `json:"noticeIds" gorm:"column:noticeIds;serializer:json"`
	NoticeRoutes          []NoticeRoute `json:"noticeRoutes" gorm:"noticeRoutes;serializer:json"`
	RepeatNoticeInterval  int64         `json:"repeatNoticeInterval"`
	RecoverNotify         *bool         `json:"recoverNotify"`
	AggregationType       string        `json:"aggregationType"`
	CreateAt              int64         `json:"createAt"`
	RecoverWaitTime       int64         `json:"recoverWaitTime"`
	CurrentPreAlertNumber int64         `json:"currentPreAlertNumber" gorm:"-"`
	CurrentAlertNumber    int64         `json:"currentAlertNumber" gorm:"-"`
	CurrentMuteNumber     int64         `json:"currentMuteNumber" gorm:"-"`
	CurrentRecoverNumber  int64         `json:"currentRecoverNumber" gorm:"-"`
}

type NoticeRoute struct {
	Key       string   `json:"key"`
	Value     string   `json:"value"`
	NoticeIds []string `json:"noticeIds" gorm:"column:noticeIds;serializer:json"`
}

func (f *FaultCenter) TableName() string {
	return "w8t_fault_center"
}

func (f *FaultCenter) GetRecoverNotify() *bool {
	if f.RecoverNotify == nil {
		isOk := false
		return &isOk
	}
	return f.RecoverNotify
}

func (f *FaultCenter) GetAlarmAggregationType() string {
	return f.AggregationType
}

func (f *FaultCenter) GetFaultCenterKey() string {
	return fmt.Sprintf("w8t:%s:%s:%s.events", f.TenantId, FaultCenterPrefix, f.ID)
}

func (f *FaultCenter) GetFaultCenterInfoKey() string {
	return fmt.Sprintf("w8t:%s:%s:%s.info", f.TenantId, FaultCenterPrefix, f.ID)
}

type FaultCenterQuery struct {
	TenantId string `form:"tenantId"`
	ID       string `form:"id"`
	Name     string `form:"name"`
	Query    string `from:"query"`
}

func BuildCacheEventKey(tenantId, faultCenterId string) string {
	return fmt.Sprintf("w8t:%s:%s:%s.events", tenantId, FaultCenterPrefix, faultCenterId)
}

func BuildCacheMuteKey(tenantId, faultCenterId string) string {
	return fmt.Sprintf("w8t:%s:%s:%s.mutes", tenantId, FaultCenterPrefix, faultCenterId)
}

func BuildCacheInfoKey(tenantId, faultCenterId string) string {
	return fmt.Sprintf("w8t:%s:%s:%s.info", tenantId, FaultCenterPrefix, faultCenterId)
}
