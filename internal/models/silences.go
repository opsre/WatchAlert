package models

type AlertSilences struct {
	TenantId      string         `json:"tenantId"`
	Name          string         `json:"name"`
	Id            string         `json:"id"`
	Labels        []SilenceLabel `json:"labels" gorm:"labels;serializer:json"`
	StartsAt      int64          `json:"startsAt"`
	UpdateBy      string         `json:"updateBy"`
	EndsAt        int64          `json:"endsAt"`
	UpdateAt      int64          `json:"updateAt"`
	FaultCenterId string         `json:"faultCenterId"`
	Comment       string         `json:"comment"`
	Status        int            `json:"status"` // 0 未生效, 1 进行中, 2 已失效
}

type SilenceLabel struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	Operator string `json:"operator"`
}

type AlertSilenceQuery struct {
	TenantId      string `json:"tenantId" form:"tenantId"`
	Id            string `json:"id" form:"id"`
	Query         string `json:"query" form:"query"`
	FaultCenterId string `json:"faultCenterId" form:"faultCenterId"`
	Status        int    `json:"status" form:"status"`
	Page
}

type SilenceResponse struct {
	List []AlertSilences `json:"list"`
	Page
}
