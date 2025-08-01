package api

import (
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
)

type dashboardController struct{}

var DashboardController = new(dashboardController)

func (dashboardController dashboardController) API(gin *gin.RouterGroup) {
	a := gin.Group("dashboard")
	a.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		a.POST("createFolder", dashboardController.CreateFolder)
		a.POST("updateFolder", dashboardController.UpdateFolder)
		a.POST("deleteFolder", dashboardController.DeleteFolder)
	}

	b := gin.Group("dashboard")
	b.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		b.GET("listFolder", dashboardController.ListFolder)
		b.GET("getFolder", dashboardController.GetFolder)
		b.GET("listGrafanaDashboards", dashboardController.ListGrafanaDashboards)
	}

	c := gin.Group("dashboard")
	c.Use(
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		c.GET("getDashboardFullUrl", dashboardController.GetDashboardFullUrl)
	}
}

func (dashboardController dashboardController) ListFolder(ctx *gin.Context) {
	r := new(types.RequestDashboardFoldersQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.ListFolder(r)
	})
}

func (dashboardController dashboardController) SearchFolder(ctx *gin.Context) {
	r := new(types.RequestDashboardFoldersQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.GetFolder(r)
	})
}

func (dashboardController dashboardController) GetFolder(ctx *gin.Context) {
	r := new(types.RequestDashboardFoldersQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.GetFolder(r)
	})
}

func (dashboardController dashboardController) CreateFolder(ctx *gin.Context) {
	r := new(types.RequestDashboardFoldersCreate)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.CreateFolder(r)
	})
}

func (dashboardController dashboardController) UpdateFolder(ctx *gin.Context) {
	r := new(types.RequestDashboardFoldersUpdate)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.UpdateFolder(r)
	})
}

func (dashboardController dashboardController) DeleteFolder(ctx *gin.Context) {
	r := new(types.RequestDashboardFoldersQuery)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.DeleteFolder(r)
	})
}

func (dashboardController dashboardController) ListGrafanaDashboards(ctx *gin.Context) {
	r := new(types.RequestDashboardFoldersQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.ListGrafanaDashboards(r)
	})
}

func (dashboardController dashboardController) GetDashboardFullUrl(ctx *gin.Context) {
	r := new(types.RequestGetGrafanaDashboard)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.GetDashboardFullUrl(r)
	})
}
