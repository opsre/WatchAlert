package api

import (
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
)

type ruleGroupController struct{}

var RuleGroupController = new(ruleGroupController)

/*
规则组 API
/api/w8t/ruleGroup
*/
func (ruleGroupController ruleGroupController) API(gin *gin.RouterGroup) {
	a := gin.Group("ruleGroup")
	a.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		a.POST("ruleGroupCreate", ruleGroupController.Create)
		a.POST("ruleGroupUpdate", ruleGroupController.Update)
		a.POST("ruleGroupDelete", ruleGroupController.Delete)
	}
	b := gin.Group("ruleGroup")
	b.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		b.GET("ruleGroupList", ruleGroupController.List)
	}
}

func (ruleGroupController ruleGroupController) Create(ctx *gin.Context) {
	r := new(types.RequestRuleGroupCreate)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RuleGroupService.Create(r)
	})
}

func (ruleGroupController ruleGroupController) Update(ctx *gin.Context) {
	r := new(types.RequestRuleGroupUpdate)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RuleGroupService.Update(r)
	})
}

func (ruleGroupController ruleGroupController) List(ctx *gin.Context) {
	r := new(types.RequestRuleGroupQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RuleGroupService.List(r)
	})
}

func (ruleGroupController ruleGroupController) Delete(ctx *gin.Context) {
	r := new(types.RequestRuleGroupQuery)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RuleGroupService.Delete(r)
	})
}
