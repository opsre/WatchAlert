package api

import (
	"errors"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
	"watchAlert/pkg/tools"

	"github.com/gin-gonic/gin"
)

type recordingRuleController struct{}

var RecordingRuleController = new(recordingRuleController)

/*
记录规则 API
/api/w8t/recordingRule
*/
func (recordingRuleController recordingRuleController) API(gin *gin.RouterGroup) {
	a := gin.Group("recordingRule")
	a.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		a.POST("recordingRuleCreate", recordingRuleController.Create)
		a.POST("recordingRuleUpdate", recordingRuleController.Update)
		a.POST("recordingRuleDelete", recordingRuleController.Delete)
	}
	b := gin.Group("recordingRule")
	b.Use(
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		b.GET("recordingRuleList", recordingRuleController.List)
		b.GET("recordingRuleGet", recordingRuleController.Get)
	}
	c := gin.Group("recordingRule")
	c.Use(
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		c.POST("recordingRuleChangeStatus", recordingRuleController.ChangeStatus)
	}
}

func (recordingRuleController recordingRuleController) Create(ctx *gin.Context) {
	r := new(types.RequestRecordingRuleCreate)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		tokenStr := ctx.Request.Header.Get("Authorization")
		if len(tokenStr) <= 0 {
			return nil, errors.New("用户未登录")
		}
		r.UpdateBy = tools.GetUser(tokenStr)

		tid, _ := ctx.Get("TenantID")
		r.TenantId = tid.(string)

		return services.RecordingRuleService.Create(r)
	})
}

func (recordingRuleController recordingRuleController) Update(ctx *gin.Context) {
	r := new(types.RequestRecordingRuleUpdate)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		tokenStr := ctx.Request.Header.Get("Authorization")
		if len(tokenStr) <= 0 {
			return nil, errors.New("用户未登录")
		}
		r.UpdateBy = tools.GetUser(tokenStr)

		tid, _ := ctx.Get("TenantID")
		r.TenantId = tid.(string)

		return services.RecordingRuleService.Update(r)
	})
}

func (recordingRuleController recordingRuleController) List(ctx *gin.Context) {
	r := new(types.RequestRecordingRuleQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RecordingRuleService.List(r)
	})
}

func (recordingRuleController recordingRuleController) Delete(ctx *gin.Context) {
	r := new(types.RequestRecordingRuleQuery)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RecordingRuleService.Delete(r)
	})
}

func (recordingRuleController recordingRuleController) Get(ctx *gin.Context) {
	r := new(types.RequestRecordingRuleQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RecordingRuleService.Get(r)
	})
}

func (recordingRuleController recordingRuleController) ChangeStatus(ctx *gin.Context) {
	r := new(types.RequestRecordingRuleChangeStatus)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RecordingRuleService.ChangeStatus(r)
	})
}
