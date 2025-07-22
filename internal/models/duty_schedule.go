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
	TenantId    string     `json:"tenantId"`
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Manager     DutyUser   `json:"manager" gorm:"manager;serializer:json"`
	Description string     `json:"description"`
	CurDutyUser []DutyUser `json:"curDutyUser" gorm:"curDutyUser;serializer:json"`
	CreateBy    string     `json:"create_by"`
	CreateAt    int64      `json:"create_at"`
}

type DutyManagementQuery struct {
	TenantId string `json:"tenantId" form:"tenantId"`
	ID       string `json:"id" form:"id"`
	Name     string `json:"name" form:"name"`
}

type DutyScheduleCreate struct {
	TenantId   string `json:"tenantId"`
	DutyId     string `json:"dutyId"`
	DutyPeriod int    `json:"dutyPeriod"`
	Month      string `json:"month"`
	//Users      []DutyUser `json:"users"`
	UserGroup [][]DutyUser `json:"userGroup"`
	DateType  string       `json:"dateType"`
	Status    string       `json:"status" `
}

type DutyUser struct {
	UserId   string `json:"userid"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
}

type CalendarStatus string

const (
	// CalendarTemporaryStatus 临时值班状态
	CalendarTemporaryStatus string = "Temporary"
	// CalendarFormalStatus 正式值班状态
	CalendarFormalStatus string = "Formal"
)

type DutySchedule struct {
	TenantId string `json:"tenantId"`
	DutyId   string `json:"dutyId"`
	Time     string `json:"time"`
	//Users    DutyUser `json:"users"`
	Status string     `json:"status"`
	Users  []DutyUser `json:"users" gorm:"users;serializer:json"`
}

type DutyScheduleQuery struct {
	TenantId string `json:"tenantId" form:"tenantId"`
	DutyId   string `json:"dutyId" form:"dutyId"`
	Time     string `json:"time" form:"time"`
	Year     string `json:"year" form:"year"`
	Month    string `json:"mouth" form:"month"`
}
