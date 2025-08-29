package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/services"
)

type FaultCenterController struct{}

func (fcc FaultCenterController) API(gin *gin.RouterGroup) {
	faultCenterA := gin.Group("faultCenter")
	faultCenterA.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		faultCenterA.POST("faultCenterCreate", fcc.Create)
		faultCenterA.POST("faultCenterUpdate", fcc.Update)
		faultCenterA.POST("faultCenterDelete", fcc.Delete)
		faultCenterA.POST("faultCenterReset", fcc.Reset)
	}

	faultCenterB := gin.Group("faultCenter")
	faultCenterB.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		faultCenterB.GET("faultCenterList", fcc.List)
		faultCenterB.GET("faultCenterSearch", fcc.Search)
	}
}

func (fcc FaultCenterController) Create(ctx *gin.Context) {
	r := new(models.FaultCenter)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.FaultCenterService.Create(r)
	})
}

func (fcc FaultCenterController) Update(ctx *gin.Context) {
	r := new(models.FaultCenter)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.FaultCenterService.Update(r)
	})
}

func (fcc FaultCenterController) Delete(ctx *gin.Context) {
	r := new(models.FaultCenterQuery)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.FaultCenterService.Delete(r)
	})
}

func (fcc FaultCenterController) List(ctx *gin.Context) {
	r := new(models.FaultCenterQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.FaultCenterService.List(r)
	})
}

func (fcc FaultCenterController) Search(ctx *gin.Context) {
	r := new(models.FaultCenterQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.FaultCenterService.Get(r)
	})
}

func (fcc FaultCenterController) Reset(ctx *gin.Context) {
	r := new(models.FaultCenter)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.FaultCenterService.Reset(r)
	})
}
