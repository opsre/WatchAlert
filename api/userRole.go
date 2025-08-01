package api

import (
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
)

type userRoleController struct{}

var UserRoleController = new(userRoleController)

/*
用户角色 API
/api/w8t/role
*/
func (userRoleController userRoleController) API(gin *gin.RouterGroup) {
	a := gin.Group("role")
	a.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		a.POST("roleCreate", userRoleController.Create)
		a.POST("roleUpdate", userRoleController.Update)
		a.POST("roleDelete", userRoleController.Delete)
	}

	b := gin.Group("role")
	b.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		b.GET("roleList", userRoleController.List)
	}
}

func (userRoleController userRoleController) Create(ctx *gin.Context) {
	r := new(types.RequestUserRoleCreate)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserRoleService.Create(r)
	})
}

func (userRoleController userRoleController) Update(ctx *gin.Context) {
	r := new(types.RequestUserRoleUpdate)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserRoleService.Update(r)
	})
}

func (userRoleController userRoleController) Delete(ctx *gin.Context) {
	r := new(types.RequestUserRoleQuery)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserRoleService.Delete(r)
	})
}

func (userRoleController userRoleController) List(ctx *gin.Context) {
	r := new(types.RequestUserRoleQuery)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserRoleService.List(r)
	})
}
