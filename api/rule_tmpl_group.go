package api

import (
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/services"
)

type RuleTmplGroupController struct{}

/*
规则模版组 API
/api/w8t/ruleTmplGroup
*/
func (rtg RuleTmplGroupController) API(gin *gin.RouterGroup) {
	ruleTmplGroupA := gin.Group("ruleTmplGroup")
	ruleTmplGroupA.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		ruleTmplGroupA.POST("ruleTmplGroupCreate", rtg.Create)
		ruleTmplGroupA.POST("ruleTmplGroupUpdate", rtg.Update)
		ruleTmplGroupA.POST("ruleTmplGroupDelete", rtg.Delete)
	}

	ruleTmplGroupB := gin.Group("ruleTmplGroup")
	ruleTmplGroupB.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		ruleTmplGroupB.GET("ruleTmplGroupList", rtg.List)
	}
}

func (rtg RuleTmplGroupController) Create(ctx *gin.Context) {
	r := new(models.RuleTemplateGroup)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RuleTmplGroupService.Create(r)
	})
}

func (rtg RuleTmplGroupController) Update(ctx *gin.Context) {
	r := new(models.RuleTemplateGroup)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RuleTmplGroupService.Update(r)
	})
}

func (rtg RuleTmplGroupController) Delete(ctx *gin.Context) {
	r := new(models.RuleTemplateGroupQuery)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RuleTmplGroupService.Delete(r)
	})
}

func (rtg RuleTmplGroupController) List(ctx *gin.Context) {
	r := new(models.RuleTemplateGroupQuery)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RuleTmplGroupService.List(r)
	})
}
