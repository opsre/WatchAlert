package api

import (
	"errors"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
	jwtUtils "watchAlert/pkg/tools"

	"github.com/gin-gonic/gin"
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

	Service(ctx, func() (interface{}, interface{}) {
		token := ctx.Request.Header.Get("Authorization")
		if token == "" {
			return nil, errors.New("token is empty")
		}

		userId := jwtUtils.GetUserID(token)
		r.UserId = userId

		return services.UserService.Info(r)
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

	Service(ctx, func() (interface{}, interface{}) {
		createUser := jwtUtils.GetUser(ctx.Request.Header.Get("Authorization"))
		if createUser == "" {
			return nil, errors.New("user is empty")
		}
		r.CreateBy = createUser

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
		return services.UserService.Check(r)
	})
}

func (userController userController) ChangePass(ctx *gin.Context) {
	r := new(types.RequestUserChangePassword)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		userID := jwtUtils.GetUserID(ctx.Request.Header.Get("Authorization"))
		if userID == "" {
			return nil, errors.New("UserID 不能为空")
		}

		if userID != r.UserId {
			return nil, errors.New("只能修改自己的密码")
		}

		return services.UserService.ChangePass(r)
	})
}
