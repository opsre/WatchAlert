package models

type Member struct {
	UserId     string   `json:"userid"`
	UserName   string   `json:"username"`
	Email      string   `json:"email"`
	Phone      string   `json:"phone"`
	Password   string   `json:"-"` // 密码哈希, 禁止序列化到响应
	Role       string   `json:"role"`
	CreateBy   string   `json:"create_by"`
	CreateAt   int64    `json:"create_at"`
	JoinDuty   string   `json:"joinDuty" `
	DutyUserId string   `json:"dutyUserId"`
	Tenants    []string `json:"tenants" gorm:"tenants;serializer:json"`
}

type ResponseLoginInfo struct {
	Token      string `json:"token"`
	Identifier string `json:"identifier"`
	UserId     string `json:"userId"`
}
