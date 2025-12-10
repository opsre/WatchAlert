package types

import "watchAlert/internal/models"

// RequestTopologyCreate 请求创建拓扑图
type RequestTopologyCreate struct {
	TenantId  string        `json:"tenantId"`
	Name      string        `json:"name"`
	Nodes     []models.Node `json:"nodes"`
	Edges     []models.Edge `json:"edges"`
	UpdatedBy string        `json:"updatedBy"`
}

// RequestTopologyUpdate 请求更新拓扑图
type RequestTopologyUpdate struct {
	TenantId  string        `json:"tenantId"`
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	Nodes     []models.Node `json:"nodes"`
	Edges     []models.Edge `json:"edges"`
	UpdatedBy string        `json:"updatedBy"`
}

// RequestTopologyQuery 请求查询拓扑图
// List接口返回简要信息（不包含nodes和edges）
// Get接口返回完整信息
type RequestTopologyQuery struct {
	TenantId string `form:"tenantId"`
	ID       string `form:"id"`
	Name     string `form:"name"`
	Query    string `form:"query"`
}

// RequestTopologyDelete 请求删除拓扑图
type RequestTopologyDelete struct {
	TenantId string `json:"tenantId"`
	ID       string `json:"id"`
}
