package api

import (
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
)

type ruleTmplGroupController struct{}

var RuleTmplGroupController = new(ruleTmplGroupController)

/*
规则模版组 API
/api/w8t/ruleTmplGroup
*/
func (ruleTmplGroupController ruleTmplGroupController) API(gin *gin.RouterGroup) {
	a := gin.Group("ruleTmplGroup")
	a.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		a.POST("ruleTmplGroupCreate", ruleTmplGroupController.Create)
		a.POST("ruleTmplGroupUpdate", ruleTmplGroupController.Update)
		a.POST("ruleTmplGroupDelete", ruleTmplGroupController.Delete)
	}

	b := gin.Group("ruleTmplGroup")
	b.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		b.GET("ruleTmplGroupList", ruleTmplGroupController.List)
	}
}

func (ruleTmplGroupController ruleTmplGroupController) Create(ctx *gin.Context) {
	r := new(types.RequestRuleTemplateGroupCreate)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RuleTmplGroupService.Create(r)
	})
}

func (ruleTmplGroupController ruleTmplGroupController) Update(ctx *gin.Context) {
	r := new(types.RequestRuleTemplateGroupUpdate)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RuleTmplGroupService.Update(r)
	})
}

func (ruleTmplGroupController ruleTmplGroupController) Delete(ctx *gin.Context) {
	r := new(types.RequestRuleTemplateGroupQuery)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RuleTmplGroupService.Delete(r)
	})
}

func (ruleTmplGroupController ruleTmplGroupController) List(ctx *gin.Context) {
	r := new(types.RequestRuleTemplateGroupQuery)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RuleTmplGroupService.List(r)
	})
}
