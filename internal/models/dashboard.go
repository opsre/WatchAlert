package models

type Dashboard struct {
	TenantId    string `json:"tenantId"`
	ID          string `json:"id" `
	Name        string `json:"name" gorm:"unique"`
	URL         string `json:"url"`
	FolderId    string `json:"folderId"`
	Description string `json:"description"`
}

type DashboardQuery struct {
	TenantId string `json:"tenantId" form:"tenantId"`
	ID       string `json:"id" form:"id"`
	Query    string `json:"query" form:"query"`
}

const (
	GrafanaV11 string = "v11"
	GrafanaV10 string = "v10"
)

type DashboardFolders struct {
	TenantId            string `json:"tenantId" form:"tenantId"`
	ID                  string `json:"id" form:"id"`
	Name                string `json:"name"`
	Theme               string `json:"theme" form:"theme"`
	GrafanaVersion      string `json:"grafanaVersion" form:"grafanaVersion"` // v10及以下, v11及以上
	GrafanaHost         string `json:"grafanaHost" form:"grafanaHost"`
	GrafanaFolderId     string `json:"grafanaFolderId" form:"grafanaFolderId"`
	GrafanaDashboardUid string `json:"grafanaDashboardUid" form:"grafanaDashboardUid" gorm:"-"`
}

type GrafanaDashboardInfo struct {
	Uid   string `json:"uid"`
	Title string `json:"title"`
	//FolderUid string `json:"folderUid"`
}

type GrafanaDashboardMeta struct {
	Meta meta `json:"meta"`
}

type meta struct {
	Url string `json:"url"`
}
