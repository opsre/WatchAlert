package models

type Member struct {
	UserId     string   `json:"userid"`
	UserName   string   `json:"username"`
	Email      string   `json:"email"`
	Phone      string   `json:"phone"`
	Password   string   `json:"password"`
	Role       string   `json:"role"`
	CreateBy   string   `json:"create_by"`
	CreateAt   int64    `json:"create_at"`
	JoinDuty   string   `json:"joinDuty" `
	DutyUserId string   `json:"dutyUserId"`
	Tenants    []string `json:"tenants" gorm:"tenants;serializer:json"`
}

type MemberQuery struct {
	UserId   string `json:"userid" form:"userid"`
	UserName string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Phone    string `json:"phone" form:"phone"`
	Query    string `json:"query" form:"query"`
	JoinDuty string `json:"joinDuty" form:"joinDuty"`
}

type ResponseLoginInfo struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	UserId   string `json:"userId"`
}
