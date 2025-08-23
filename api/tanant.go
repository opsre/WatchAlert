package api

import (
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
	jwtUtils "watchAlert/pkg/tools"
)

type tenantController struct{}

var TenantController = new(tenantController)

/*
租户 API
/api/w8t/tenant
*/
func (tenantController tenantController) API(gin *gin.RouterGroup) {
	a := gin.Group("tenant")
	a.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.AuditingLog(),
	)
	{
		a.POST("createTenant", tenantController.Create)
		a.POST("updateTenant", tenantController.Update)
		a.POST("deleteTenant", tenantController.Delete)
		a.POST("addUsersToTenant", tenantController.AddUsersToTenant)
		a.POST("delUsersOfTenant", tenantController.DelUsersOfTenant)
		a.POST("changeTenantUserRole", tenantController.ChangeTenantUserRole)
	}

	b := gin.Group("tenant")
	b.Use(
		middleware.Auth(),
		middleware.Permission(),
	)
	{
		b.GET("getTenantList", tenantController.List)
		b.GET("getTenant", tenantController.Get)
		b.GET("getUsersForTenant", tenantController.GetUsersForTenant)
	}
}

func (tenantController tenantController) Create(ctx *gin.Context) {
	r := new(types.RequestTenantCreate)
	BindJson(ctx, r)

	token := ctx.Request.Header.Get("Authorization")
	r.CreateBy = jwtUtils.GetUser(token)
	r.UserId = jwtUtils.GetUserID(token)
	if r.UserId == "" {
		r.UserId = "admin"
	}

	Service(ctx, func() (interface{}, interface{}) {
		return services.TenantService.Create(r)
	})
}

func (tenantController tenantController) Update(ctx *gin.Context) {
	r := new(types.RequestTenantUpdate)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TenantService.Update(r)
	})
}

func (tenantController tenantController) Delete(ctx *gin.Context) {
	r := new(types.RequestTenantQuery)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TenantService.Delete(r)
	})
}

func (tenantController tenantController) List(ctx *gin.Context) {
	r := new(types.RequestTenantQuery)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TenantService.List(r)
	})
}

func (tenantController tenantController) Get(ctx *gin.Context) {
	r := new(types.RequestTenantQuery)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TenantService.Get(r)
	})
}

func (tenantController tenantController) Search(ctx *gin.Context) {
	// TODO
}

func (tenantController tenantController) AddUsersToTenant(ctx *gin.Context) {
	r := new(types.RequestTenantAddUsers)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TenantService.AddUsersToTenant(r)
	})
}

func (tenantController tenantController) DelUsersOfTenant(ctx *gin.Context) {
	r := new(types.RequestTenantQuery)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TenantService.DelUsersOfTenant(r)
	})
}

func (tenantController tenantController) GetUsersForTenant(ctx *gin.Context) {
	r := new(types.RequestTenantQuery)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TenantService.GetUsersForTenant(r)
	})
}

func (tenantController tenantController) ChangeTenantUserRole(ctx *gin.Context) {
	r := new(types.RequestTenantChangeUserRole)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TenantService.ChangeTenantUserRole(r)
	})
}
