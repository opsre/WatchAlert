package api

import (
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"

	"github.com/gin-gonic/gin"
)

type recordingRuleGroupController struct{}

var RecordingRuleGroupController = new(recordingRuleGroupController)

/*
记录规则组 API
/api/w8t/recordingRuleGroup
*/
func (recordingRuleGroupController recordingRuleGroupController) API(gin *gin.RouterGroup) {
	a := gin.Group("recordingRuleGroup")
	a.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		a.POST("recordingRuleGroupCreate", recordingRuleGroupController.Create)
		a.POST("recordingRuleGroupUpdate", recordingRuleGroupController.Update)
		a.POST("recordingRuleGroupDelete", recordingRuleGroupController.Delete)
	}
	b := gin.Group("recordingRuleGroup")
	b.Use(
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		b.GET("recordingRuleGroupList", recordingRuleGroupController.List)
		b.GET("recordingRuleGroupGet", recordingRuleGroupController.Get)
	}
}

func (recordingRuleGroupController recordingRuleGroupController) Create(ctx *gin.Context) {
	r := new(types.RequestRecordingRuleGroupCreate)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RecordingRuleGroupService.Create(r)
	})
}

func (recordingRuleGroupController recordingRuleGroupController) Update(ctx *gin.Context) {
	r := new(types.RequestRecordingRuleGroupUpdate)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RecordingRuleGroupService.Update(r)
	})
}

func (recordingRuleGroupController recordingRuleGroupController) List(ctx *gin.Context) {
	r := new(types.RequestRecordingRuleGroupQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RecordingRuleGroupService.List(r)
	})
}

func (recordingRuleGroupController recordingRuleGroupController) Delete(ctx *gin.Context) {
	r := new(types.RequestRecordingRuleGroupQuery)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RecordingRuleGroupService.Delete(r)
	})
}

func (recordingRuleGroupController recordingRuleGroupController) Get(ctx *gin.Context) {
	r := new(types.RequestRecordingRuleGroupQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RecordingRuleGroupService.Get(r)
	})
}
