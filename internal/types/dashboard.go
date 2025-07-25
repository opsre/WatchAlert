package types

const (
	GrafanaV11 string = "v11"
	GrafanaV10 string = "v10"
)

type RequestDashboardQuery struct {
	TenantId string `json:"tenantId" form:"tenantId"`
	ID       string `json:"id" form:"id"`
	Query    string `json:"query" form:"query"`
}

type RequestDashboardCreate struct {
	TenantId    string `json:"tenantId"`
	Name        string `json:"name" gorm:"unique"`
	URL         string `json:"url"`
	FolderId    string `json:"folderId"`
	Description string `json:"description"`
}

type RequestDashboardUpdate struct {
	TenantId    string `json:"tenantId"`
	ID          string `json:"id" `
	Name        string `json:"name" gorm:"unique"`
	URL         string `json:"url"`
	FolderId    string `json:"folderId"`
	Description string `json:"description"`
}

type RequestDashboardFoldersQuery struct {
	TenantId string `json:"tenantId" form:"tenantId"`
	ID       string `json:"id" form:"id"`
	Query    string `json:"query" form:"query"`
}

type RequestDashboardFoldersCreate struct {
	TenantId            string `json:"tenantId" form:"tenantId"`
	Name                string `json:"name"`
	Theme               string `json:"theme" form:"theme"`
	GrafanaVersion      string `json:"grafanaVersion" form:"grafanaVersion"` // v10及以下, v11及以上
	GrafanaHost         string `json:"grafanaHost" form:"grafanaHost"`
	GrafanaFolderId     string `json:"grafanaFolderId" form:"grafanaFolderId"`
	GrafanaDashboardUid string `json:"grafanaDashboardUid" form:"grafanaDashboardUid" gorm:"-"`
}

type RequestDashboardFoldersUpdate struct {
	TenantId            string `json:"tenantId" form:"tenantId"`
	ID                  string `json:"id" form:"id"`
	Name                string `json:"name"`
	Theme               string `json:"theme" form:"theme"`
	GrafanaVersion      string `json:"grafanaVersion" form:"grafanaVersion"` // v10及以下, v11及以上
	GrafanaHost         string `json:"grafanaHost" form:"grafanaHost"`
	GrafanaFolderId     string `json:"grafanaFolderId" form:"grafanaFolderId"`
	GrafanaDashboardUid string `json:"grafanaDashboardUid" form:"grafanaDashboardUid" gorm:"-"`
}

type RequestGetGrafanaDashboard struct {
	Theme string `json:"theme" form:"theme"`
	Host  string `json:"host" form:"host"`
	Uid   string `json:"uid" form:"uid"`
}

type ResponseGrafanaDashboardInfo struct {
	Uid   string `json:"uid"`
	Title string `json:"title"`
}

type ResponseGrafanaDashboardMeta struct {
	Meta meta `json:"meta"`
}

type meta struct {
	Url string `json:"url"`
}
