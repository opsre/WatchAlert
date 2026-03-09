package api

import (
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
	jwtUtils "watchAlert/pkg/tools"

	"github.com/gin-gonic/gin"
)

type dutyController struct{}

var DutyController = new(dutyController)

/*
排班管理 API
/api/w8t/dutyManage
*/
func (dutyController dutyController) API(gin *gin.RouterGroup) {
	a := gin.Group("dutyManage")
	a.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		a.POST("dutyManageCreate", dutyController.Create)
		a.POST("dutyManageUpdate", dutyController.Update)
		a.POST("dutyManageDelete", dutyController.Delete)
	}

	b := gin.Group("dutyManage")
	b.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		b.GET("dutyManageList", dutyController.List)
	}
}

func (dutyController dutyController) List(ctx *gin.Context) {
	r := new(types.RequestDutyManagementQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DutyManageService.List(r)
	})
}

func (dutyController dutyController) Create(ctx *gin.Context) {
	r := new(types.RequestDutyManagementCreate)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		userName := jwtUtils.GetUser(ctx.Request.Header.Get("Authorization"))
		r.UpdateBy = userName

		tid, _ := ctx.Get("TenantID")
		r.TenantId = tid.(string)

		return services.DutyManageService.Create(r)
	})
}

func (dutyController dutyController) Update(ctx *gin.Context) {
	r := new(types.RequestDutyManagementUpdate)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		userName := jwtUtils.GetUser(ctx.Request.Header.Get("Authorization"))
		r.UpdateBy = userName

		tid, _ := ctx.Get("TenantID")
		r.TenantId = tid.(string)

		return services.DutyManageService.Update(r)
	})
}

func (dutyController dutyController) Delete(ctx *gin.Context) {
	r := new(types.RequestDutyManagementQuery)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DutyManageService.Delete(r)
	})
}

func (dutyController dutyController) Get(ctx *gin.Context) {
	r := new(types.RequestDutyManagementQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DutyManageService.Get(r)
	})
}
