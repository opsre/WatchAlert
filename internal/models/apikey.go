package models

import "time"

type ApiKey struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId      string    `json:"userId" gorm:"column:user_id;not null"`
	Name        string    `json:"name" gorm:"column:name;size:255;not null"`
	Description string    `json:"description" gorm:"column:description;size:500"`
	Key         string    `json:"key" gorm:"column:key;size:255;not null;uniqueIndex"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (ApiKey) TableName() string {
	return "w8t_api_keys"
}
