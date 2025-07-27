package api

import (
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
)

type DashboardController struct{}

func (dc DashboardController) API(gin *gin.RouterGroup) {
	dashboardA := gin.Group("dashboard")
	dashboardA.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		dashboardA.POST("createFolder", dc.CreateFolder)
		dashboardA.POST("updateFolder", dc.UpdateFolder)
		dashboardA.POST("deleteFolder", dc.DeleteFolder)
	}

	dashboardB := gin.Group("dashboard")
	dashboardB.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		dashboardB.GET("listFolder", dc.ListFolder)
		dashboardB.GET("getFolder", dc.GetFolder)
		dashboardB.GET("listGrafanaDashboards", dc.ListGrafanaDashboards)
	}

	c := gin.Group("dashboard")
	c.Use(
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		c.GET("getDashboardFullUrl", dc.GetDashboardFullUrl)
	}
}

func (dc DashboardController) ListFolder(ctx *gin.Context) {
	r := new(types.RequestDashboardFoldersQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.ListFolder(r)
	})
}

func (dc DashboardController) SearchFolder(ctx *gin.Context) {
	r := new(types.RequestDashboardFoldersQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.GetFolder(r)
	})
}

func (dc DashboardController) GetFolder(ctx *gin.Context) {
	r := new(types.RequestDashboardFoldersQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.GetFolder(r)
	})
}

func (dc DashboardController) CreateFolder(ctx *gin.Context) {
	r := new(types.RequestDashboardFoldersCreate)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.CreateFolder(r)
	})
}

func (dc DashboardController) UpdateFolder(ctx *gin.Context) {
	r := new(types.RequestDashboardFoldersUpdate)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.UpdateFolder(r)
	})
}

func (dc DashboardController) DeleteFolder(ctx *gin.Context) {
	r := new(types.RequestDashboardFoldersQuery)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.DeleteFolder(r)
	})
}

func (dc DashboardController) ListGrafanaDashboards(ctx *gin.Context) {
	r := new(types.RequestDashboardFoldersQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.ListGrafanaDashboards(r)
	})
}

func (dc DashboardController) GetDashboardFullUrl(ctx *gin.Context) {
	r := new(types.RequestGetGrafanaDashboard)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.GetDashboardFullUrl(r)
	})
}
