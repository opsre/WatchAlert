package api

import (
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/services"
)

type ProbingController struct{}

func (e ProbingController) API(gin *gin.RouterGroup) {
	eventA := gin.Group("probing")
	eventA.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		eventA.POST("createProbing", e.Create)
		eventA.POST("updateProbing", e.Update)
		eventA.POST("deleteProbing", e.Delete)
		eventA.POST("onceProbing", e.Once)
	}

	eventB := gin.Group("probing")
	eventB.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		eventB.GET("listProbing", e.List)
		eventB.GET("searchProbing", e.Search)
		eventB.GET("getProbingHistory", e.GetHistory)
	}

	c := gin.Group("probing")
	c.Use(
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		c.POST("changeState", e.ChangeState)
	}
}

func (e ProbingController) List(ctx *gin.Context) {
	r := new(models.ProbingRuleQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.ProbingService.List(r)
	})
}

func (e ProbingController) Search(ctx *gin.Context) {
	r := new(models.ProbingRuleQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.ProbingService.Search(r)
	})
}

func (e ProbingController) Create(ctx *gin.Context) {
	r := new(models.ProbingRule)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.ProbingService.Create(r)
	})
}

func (e ProbingController) Update(ctx *gin.Context) {
	r := new(models.ProbingRule)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.ProbingService.Update(r)
	})
}

func (e ProbingController) Delete(ctx *gin.Context) {
	r := new(models.ProbingRuleQuery)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.ProbingService.Delete(r)
	})
}

func (e ProbingController) Once(ctx *gin.Context) {
	r := new(models.OnceProbing)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.ProbingService.Once(r)
	})
}

func (e ProbingController) GetHistory(ctx *gin.Context) {
	r := new(models.ReqProbingHistory)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.ProbingService.GetHistory(r)
	})
}

func (e ProbingController) ChangeState(ctx *gin.Context) {
	r := new(models.RequestProbeChangeState)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.ProbingService.ChangeState(r)
	})
}
