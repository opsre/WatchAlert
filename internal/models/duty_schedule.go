package models

type DutyManagement struct {
	TenantId    string     `json:"tenantId"`
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Manager     DutyUser   `json:"manager" gorm:"manager;serializer:json"`
	Description string     `json:"description"`
	CurDutyUser []DutyUser `json:"curDutyUser" gorm:"curDutyUser;serializer:json"`
	UpdateBy    string     `json:"updateBy"`
	UpdateAt    int64      `json:"updateAt"`
}

type CalendarStatus string

const (
	// CalendarTemporaryStatus 临时值班状态
	CalendarTemporaryStatus string = "Temporary"
	// CalendarFormalStatus 正式值班状态
	CalendarFormalStatus string = "Formal"
)

type DutySchedule struct {
	TenantId string     `json:"tenantId"`
	DutyId   string     `json:"dutyId"`
	Time     string     `json:"time"`
	Status   string     `json:"status"`
	Users    []DutyUser `json:"users" gorm:"users;serializer:json"`
}

type DutyUser struct {
	UserId   string `json:"userid"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
}
