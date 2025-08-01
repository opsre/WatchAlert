package types

import "watchAlert/internal/models"

// RequestUserRoleCreate 请求创建用户角色
type RequestUserRoleCreate struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Permissions []models.UserPermissions `json:"permissions" gorm:"permissions;serializer:json"`
	CreateAt    int64                    `json:"create_at"`
}

// RequestUserRoleUpdate 请求更新用户角色
type RequestUserRoleUpdate struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Permissions []models.UserPermissions `json:"permissions" gorm:"permissions;serializer:json"`
	CreateAt    int64                    `json:"create_at"`
}

// RequestUserRoleQuery 请求查询用户角色
type RequestUserRoleQuery struct {
	ID          string `json:"id" form:"id"`
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
}
