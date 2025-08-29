package api

import (
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/services"
)

type RuleTmplController struct{}

/*
规则模版 API
/api/w8t/ruleTmpl
*/
func (rtc RuleTmplController) API(gin *gin.RouterGroup) {
	ruleTmplA := gin.Group("ruleTmpl")
	ruleTmplA.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		ruleTmplA.POST("ruleTmplCreate", rtc.Create)
		ruleTmplA.POST("ruleTmplUpdate", rtc.Update)
		ruleTmplA.POST("ruleTmplDelete", rtc.Delete)
	}

	ruleTmplB := gin.Group("ruleTmpl")
	ruleTmplB.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		ruleTmplB.GET("ruleTmplList", rtc.List)
	}
}

func (rtc RuleTmplController) Create(ctx *gin.Context) {
	r := new(models.RuleTemplate)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RuleTmplService.Create(r)
	})
}

func (rtc RuleTmplController) Update(ctx *gin.Context) {
	r := new(models.RuleTemplate)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RuleTmplService.Update(r)
	})
}

func (rtc RuleTmplController) Delete(ctx *gin.Context) {
	r := new(models.RuleTemplateQuery)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RuleTmplService.Delete(r)
	})
}

func (rtc RuleTmplController) List(ctx *gin.Context) {
	r := new(models.RuleTemplateQuery)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RuleTmplService.List(r)
	})
}
