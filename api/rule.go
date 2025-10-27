package api

import (
	"errors"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
	"watchAlert/pkg/tools"

	"github.com/gin-gonic/gin"
)

type ruleController struct{}

var RuleController = new(ruleController)

/*
告警规则 API
/api/w8t/rule
*/
func (ruleController ruleController) API(gin *gin.RouterGroup) {
	a := gin.Group("rule")
	a.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		a.POST("ruleCreate", ruleController.Create)
		a.POST("ruleUpdate", ruleController.Update)
		a.POST("ruleDelete", ruleController.Delete)
	}
	b := gin.Group("rule")
	b.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		b.GET("ruleList", ruleController.List)
		b.GET("ruleSearch", ruleController.Search)
	}
	c := gin.Group("rule")
	c.Use(
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		c.POST("import", ruleController.Import)
		c.POST("ruleChangeStatus", ruleController.ChangeStatus)
	}
}

func (ruleController ruleController) Create(ctx *gin.Context) {
	r := new(types.RequestRuleCreate)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		tokenStr := ctx.Request.Header.Get("Authorization")
		if len(tokenStr) <= 0 {
			return nil, errors.New("用户未登录")
		}
		r.UpdateBy = tools.GetUser(tokenStr)

		tid, _ := ctx.Get("TenantID")
		r.TenantId = tid.(string)

		return services.RuleService.Create(r)
	})
}

func (ruleController ruleController) Update(ctx *gin.Context) {
	r := new(types.RequestRuleUpdate)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		tokenStr := ctx.Request.Header.Get("Authorization")
		if len(tokenStr) <= 0 {
			return nil, errors.New("用户未登录")
		}
		r.UpdateBy = tools.GetUser(tokenStr)

		tid, _ := ctx.Get("TenantID")
		r.TenantId = tid.(string)

		return services.RuleService.Update(r)
	})
}

func (ruleController ruleController) List(ctx *gin.Context) {
	r := new(types.RequestRuleQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RuleService.List(r)
	})
}

func (ruleController ruleController) Delete(ctx *gin.Context) {
	r := new(types.RequestRuleQuery)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RuleService.Delete(r)
	})
}

func (ruleController ruleController) Search(ctx *gin.Context) {
	r := new(types.RequestRuleQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RuleService.Get(r)
	})
}

func (ruleController ruleController) ChangeStatus(ctx *gin.Context) {
	r := new(types.RequestRuleChangeStatus)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RuleService.ChangeStatus(r)
	})
}

func (ruleController ruleController) Import(ctx *gin.Context) {
	r := new(types.RequestRuleImport)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.RuleService.Import(r)
	})
}
