package api

import (
	"errors"
	"watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
	"watchAlert/pkg/tools"

	"github.com/gin-gonic/gin"
)

type topologyController struct{}

var TopologyController = new(topologyController)

func (topologyController topologyController) API(gin *gin.RouterGroup) {
	a := gin.Group("topology")
	a.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		a.POST("create", topologyController.Create)
		a.POST("update", topologyController.Update)
		a.POST("delete", topologyController.Delete)
	}

	b := gin.Group("topology")
	b.Use(
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		b.GET("list", topologyController.List)
		b.GET("getDetail", topologyController.GetDetail)
	}
}

func (topologyController topologyController) Create(ctx *gin.Context) {
	r := new(types.RequestTopologyCreate)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		tokenStr := ctx.Request.Header.Get("Authorization")
		if len(tokenStr) <= 0 {
			return nil, errors.New("用户未登录")
		}
		r.UpdatedBy = tools.GetUser(tokenStr)

		return services.TopologyService.Create(r)
	})
}

func (topologyController topologyController) Update(ctx *gin.Context) {
	r := new(types.RequestTopologyUpdate)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		tokenStr := ctx.Request.Header.Get("Authorization")
		if len(tokenStr) <= 0 {
			return nil, errors.New("用户未登录")
		}
		r.UpdatedBy = tools.GetUser(tokenStr)

		return services.TopologyService.Update(r)
	})
}

func (topologyController topologyController) Delete(ctx *gin.Context) {
	r := new(types.RequestTopologyDelete)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TopologyService.Delete(r)
	})
}

func (topologyController topologyController) List(ctx *gin.Context) {
	r := new(types.RequestTopologyQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TopologyService.List(r)
	})
}

func (topologyController topologyController) GetDetail(ctx *gin.Context) {
	r := new(types.RequestTopologyQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TopologyService.GetDetail(r)
	})
}
