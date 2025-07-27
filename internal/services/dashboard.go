package services

import (
	"fmt"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
	"watchAlert/pkg/tools"
)

type dashboardService struct {
	ctx *ctx.Context
}

type InterDashboardService interface {
	ListFolder(req interface{}) (data interface{}, error interface{})
	GetFolder(req interface{}) (data interface{}, error interface{})
	CreateFolder(req interface{}) (data interface{}, error interface{})
	UpdateFolder(req interface{}) (data interface{}, error interface{})
	DeleteFolder(req interface{}) (data interface{}, error interface{})
	ListGrafanaDashboards(req interface{}) (data interface{}, error interface{})
	GetDashboardFullUrl(req interface{}) (data interface{}, error interface{})
}

func newInterDashboardService(ctx *ctx.Context) InterDashboardService {
	return &dashboardService{
		ctx: ctx,
	}
}

func (ds dashboardService) ListFolder(req interface{}) (data interface{}, error interface{}) {
	r := req.(*types.RequestDashboardFoldersQuery)
	folder, err := ds.ctx.DB.Dashboard().ListDashboardFolder(r.TenantId, r.Query)
	if err != nil {
		return nil, err
	}

	return folder, nil
}

func (ds dashboardService) GetFolder(req interface{}) (data interface{}, error interface{}) {
	r := req.(*types.RequestDashboardFoldersQuery)

	folder, err := ds.ctx.DB.Dashboard().GetDashboardFolder(r.TenantId, r.ID)
	if err != nil {
		return nil, err
	}

	return folder, nil
}

func (ds dashboardService) CreateFolder(req interface{}) (data interface{}, error interface{}) {
	r := req.(*types.RequestDashboardFoldersCreate)
	err := ctx.DB.Dashboard().CreateDashboardFolder(models.DashboardFolders{
		TenantId:            r.TenantId,
		ID:                  "f-" + tools.RandId(),
		Name:                r.Name,
		Theme:               r.Theme,
		GrafanaVersion:      r.GrafanaVersion,
		GrafanaHost:         r.GrafanaHost,
		GrafanaFolderId:     r.GrafanaFolderId,
		GrafanaDashboardUid: r.GrafanaDashboardUid,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ds dashboardService) UpdateFolder(req interface{}) (data interface{}, error interface{}) {
	r := req.(*types.RequestDashboardFoldersUpdate)
	err := ctx.DB.Dashboard().UpdateDashboardFolder(models.DashboardFolders{
		TenantId:            r.TenantId,
		ID:                  r.ID,
		Name:                r.Name,
		Theme:               r.Theme,
		GrafanaVersion:      r.GrafanaVersion,
		GrafanaHost:         r.GrafanaHost,
		GrafanaFolderId:     r.GrafanaFolderId,
		GrafanaDashboardUid: r.GrafanaDashboardUid,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ds dashboardService) DeleteFolder(req interface{}) (data interface{}, error interface{}) {
	r := req.(*types.RequestDashboardFoldersQuery)
	err := ctx.DB.Dashboard().DeleteDashboardFolder(r.TenantId, r.ID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ds dashboardService) ListGrafanaDashboards(req interface{}) (data interface{}, error interface{}) {
	r := req.(*types.RequestDashboardFoldersQuery)

	folder, err := ds.ctx.DB.Dashboard().GetDashboardFolder(r.TenantId, r.ID)
	if err != nil {
		return nil, err
	}

	var query string
	switch folder.GrafanaVersion {
	case types.GrafanaV11:
		query = fmt.Sprintf("folderUIDs=%s&deleted=false&limit=1000", folder.GrafanaFolderId)
	case types.GrafanaV10:
		query = fmt.Sprintf("folderIds=%s", folder.GrafanaFolderId)
	default:
		return nil, fmt.Errorf("invalid grafana version, please change v10 or v11")
	}

	requestURL := fmt.Sprintf("%s/api/search?%s", folder.GrafanaHost, query)
	get, err := tools.Get(nil, requestURL, 10)
	if err != nil {
		return nil, fmt.Errorf("请求错误, err: %s", err.Error())
	}

	var d []types.ResponseGrafanaDashboardInfo
	if err := tools.ParseReaderBody(get.Body, &d); err != nil {
		return nil, fmt.Errorf("读取body错误, err: %s", err.Error())
	}

	return d, nil
}

func (ds dashboardService) GetDashboardFullUrl(req interface{}) (data interface{}, error interface{}) {
	r := req.(*types.RequestGetGrafanaDashboard)
	get, err := tools.Get(nil, fmt.Sprintf("%s/api/dashboards/uid/%s", r.Host, r.Uid), 10)
	if err != nil {
		return nil, err
	}

	var d types.ResponseGrafanaDashboardMeta
	if err := tools.ParseReaderBody(get.Body, &d); err != nil {
		return nil, err
	}

	full := r.Host + d.Meta.Url + "?theme=" + r.Theme
	return full, nil
}
