package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
)

type subscribeController struct{}

var SubscribeController = new(subscribeController)

func (subscribeController subscribeController) API(gin *gin.RouterGroup) {
	a := gin.Group("subscribe")
	a.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.AuditingLog(),
		middleware.ParseTenant(),
	)
	{
		a.POST("createSubscribe", subscribeController.Create)
		a.POST("deleteSubscribe", subscribeController.Delete)
	}

	b := gin.Group("subscribe")
	b.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		b.GET("listSubscribe", subscribeController.List)
		b.GET("getSubscribe", subscribeController.Get)
	}
}

func (subscribeController subscribeController) List(ctx *gin.Context) {
	r := new(types.RequestSubscribeQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.STenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.SubscribeService.List(r)
	})
}

func (subscribeController subscribeController) Get(ctx *gin.Context) {
	r := new(types.RequestSubscribeQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.STenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.SubscribeService.Get(r)
	})
}

func (subscribeController subscribeController) Create(ctx *gin.Context) {
	r := new(types.RequestSubscribeCreate)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.STenantId = tid.(string)
	uid, _ := ctx.Get("UserId")
	r.SUserId = uid.(string)
	ue, _ := ctx.Get("UserEmail")
	r.SUserEmail = ue.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.SubscribeService.Create(r)
	})
}

func (subscribeController subscribeController) Delete(ctx *gin.Context) {
	r := new(types.RequestSubscribeQuery)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.STenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.SubscribeService.Delete(r)
	})
}
