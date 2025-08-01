package api

import (
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
	jwtUtils "watchAlert/pkg/tools"
)

type silenceController struct{}

var SilenceController = new(silenceController)

/*
告警静默 API
/api/w8t/silence
*/
func (silenceController silenceController) API(gin *gin.RouterGroup) {
	a := gin.Group("silence")
	a.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		a.POST("silenceCreate", silenceController.Create)
		a.POST("silenceUpdate", silenceController.Update)
		a.POST("silenceDelete", silenceController.Delete)
	}

	b := gin.Group("silence")
	b.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		b.GET("silenceList", silenceController.List)
	}
}

func (silenceController silenceController) Create(ctx *gin.Context) {
	r := new(types.RequestSilenceCreate)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	user := jwtUtils.GetUser(ctx.Request.Header.Get("Authorization"))
	r.UpdateBy = user

	Service(ctx, func() (interface{}, interface{}) {
		return services.SilenceService.Create(r)
	})
}

func (silenceController silenceController) Update(ctx *gin.Context) {
	r := new(types.RequestSilenceUpdate)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	user := jwtUtils.GetUser(ctx.Request.Header.Get("Authorization"))
	r.UpdateBy = user

	Service(ctx, func() (interface{}, interface{}) {
		return services.SilenceService.Update(r)
	})
}

func (silenceController silenceController) Delete(ctx *gin.Context) {
	r := new(types.RequestSilenceQuery)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.SilenceService.Delete(r)
	})
}

func (silenceController silenceController) List(ctx *gin.Context) {
	r := new(types.RequestSilenceQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.SilenceService.List(r)
	})
}
