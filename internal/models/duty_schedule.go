package models

import (
	"gorm.io/gorm"
)

type People struct {
	gorm.Model
	UserID         string `gorm:"column:userId" json:"userId"`
	UserName       string `gorm:"column:userName" json:"userName"`
	Phone          string `gorm:"column:phone" json:"phone"`
	Email          string `gorm:"column:email" json:"email"`
	Notice         string `gorm:"column:sender" json:"sender"`
	FeiShuUserID   string `gorm:"column:feiShuUserID" json:"feiShuUserID"`
	DingDingUserID string `gorm:"column:dingDingUserID" json:"dingDingUserID"`
}

type PeopleGroup struct {
	gorm.Model
	GroupID   uint   `gorm:"column:groupID" json:"groupID"`
	GroupName string `gorm:"column:groupName" json:"groupName"`
}

type JoinsPeopleGroup struct {
	UserName  string
	GroupName string
}

type DutyManagement struct {
	TenantId    string `json:"tenantId"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Manager     Users  `json:"manager" gorm:"manager;serializer:json"`
	Description string `json:"description"`
	CurDutyUser string `json:"curDutyUser"`
	CreateBy    string `json:"create_by"`
	CreateAt    int64  `json:"create_at"`
}

type DutyManagementQuery struct {
	TenantId string `json:"tenantId" form:"tenantId"`
	ID       string `json:"id" form:"id"`
	Name     string `json:"name" form:"name"`
}

type DutyScheduleCreate struct {
	TenantId   string  `json:"tenantId"`
	DutyId     string  `json:"dutyId"`
	DutyPeriod int     `json:"dutyPeriod"`
	Month      string  `json:"month"`
	Users      []Users `json:"users"`
	DateType   string  `json:"dateType"`
}

type Users struct {
	UserId   string `json:"userid"`
	Username string `json:"username"`
}

type DutySchedule struct {
	TenantId string `json:"tenantId"`
	DutyId   string `json:"dutyId"`
	Time     string `json:"time"`
	Users
}

type DutyScheduleQuery struct {
	TenantId string `json:"tenantId" form:"tenantId"`
	DutyId   string `json:"dutyId" form:"dutyId"`
	Time     string `json:"time" form:"time"`
	Year     string `json:"year" form:"year"`
	Month    string `json:"mouth" form:"month"`
}
