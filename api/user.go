package api

import (
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
	jwtUtils "watchAlert/pkg/tools"
)

type userController struct{}

var UserController = new(userController)

/*
用户 API
/api/w8t/user
*/
func (userController userController) API(gin *gin.RouterGroup) {

	a := gin.Group("user")
	a.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		a.POST("userUpdate", userController.Update)
		a.POST("userDelete", userController.Delete)
		a.POST("userChangePass", userController.ChangePass)
	}

	b := gin.Group("user")
	b.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		b.GET("userList", userController.List)
	}

}

func (userController userController) List(ctx *gin.Context) {
	r := new(types.RequestUserQuery)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.List(r)
	})
}

func (userController userController) GetUserInfo(ctx *gin.Context) {
	r := new(types.RequestUserQuery)
	BindQuery(ctx, r)

	token := ctx.Request.Header.Get("Authorization")
	username := jwtUtils.GetUser(token)
	r.UserName = username

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.Get(r)
	})
}

func (userController userController) Login(ctx *gin.Context) {
	r := new(types.RequestUserLogin)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.Login(r)
	})
}

func (userController userController) Register(ctx *gin.Context) {
	r := new(types.RequestUserCreate)
	BindJson(ctx, r)

	createUser := jwtUtils.GetUser(ctx.Request.Header.Get("Authorization"))
	r.CreateBy = createUser

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.Register(r)
	})
}

func (userController userController) Update(ctx *gin.Context) {
	r := new(types.RequestUserUpdate)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.Update(r)
	})
}

func (userController userController) Delete(ctx *gin.Context) {
	r := new(types.RequestUserQuery)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.Delete(r)
	})
}

func (userController userController) CheckUser(ctx *gin.Context) {
	r := new(types.RequestUserQuery)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.Get(r)
	})
}

func (userController userController) ChangePass(ctx *gin.Context) {
	r := new(types.RequestUserChangePassword)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.ChangePass(r)
	})
}
