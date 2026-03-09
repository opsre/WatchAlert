package api

import (
	"watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"

	"github.com/gin-gonic/gin"
)

type faultCenterController struct{}

var FaultCenterController = new(faultCenterController)

func (faultCenterController faultCenterController) API(gin *gin.RouterGroup) {
	faultCenterA := gin.Group("faultCenter")
	faultCenterA.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		faultCenterA.POST("faultCenterCreate", faultCenterController.Create)
		faultCenterA.POST("faultCenterUpdate", faultCenterController.Update)
		faultCenterA.POST("faultCenterDelete", faultCenterController.Delete)
		faultCenterA.POST("faultCenterReset", faultCenterController.Reset)
	}

	faultCenterB := gin.Group("faultCenter")
	faultCenterB.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		faultCenterB.GET("faultCenterList", faultCenterController.List)
		faultCenterB.GET("faultCenterSearch", faultCenterController.Search)
	}

	c := gin.Group("faultCenter")
	c.Use(
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		c.GET("slo", faultCenterController.Slo)
	}
}

func (faultCenterController faultCenterController) Create(ctx *gin.Context) {
	r := new(types.RequestFaultCenterCreate)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.FaultCenterService.Create(r)
	})
}

func (faultCenterController faultCenterController) Update(ctx *gin.Context) {
	r := new(types.RequestFaultCenterUpdate)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.FaultCenterService.Update(r)
	})
}

func (faultCenterController faultCenterController) Delete(ctx *gin.Context) {
	r := new(types.RequestFaultCenterQuery)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.FaultCenterService.Delete(r)
	})
}

func (faultCenterController faultCenterController) List(ctx *gin.Context) {
	r := new(types.RequestFaultCenterQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.FaultCenterService.List(r)
	})
}

func (faultCenterController faultCenterController) Search(ctx *gin.Context) {
	r := new(types.RequestFaultCenterQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.FaultCenterService.Get(r)
	})
}

func (faultCenterController faultCenterController) Reset(ctx *gin.Context) {
	r := new(types.RequestFaultCenterReset)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.FaultCenterService.Reset(r)
	})
}

func (faultCenterController faultCenterController) Slo(ctx *gin.Context) {
	r := new(types.RequestFaultCenterQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.FaultCenterService.Slo(r)
	})
}
