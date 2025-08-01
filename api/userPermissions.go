package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/internal/middleware"
	"watchAlert/internal/services"
)

type userPermissionsController struct{}

var UserPermissionsController = new(userPermissionsController)

/*
用户权限 API
/api/w8t/permissions
*/
func (userPermissionsController userPermissionsController) API(gin *gin.RouterGroup) {
	a := gin.Group("permissions")
	a.Use(
		middleware.Auth(),
	)
	{
		a.GET("permsList", userPermissionsController.List)
	}
}

func (userPermissionsController userPermissionsController) List(ctx *gin.Context) {
	Service(ctx, func() (interface{}, interface{}) {
		return services.UserPermissionService.List()
	})
}
