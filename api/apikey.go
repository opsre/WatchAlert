package api

import (
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
	"watchAlert/pkg/tools"

	"github.com/gin-gonic/gin"
)

type apiKeyController struct{}

var ApiKeyController = new(apiKeyController)

/*
API Key API
/api/w8t/apikey
*/
func (apiKeyController apiKeyController) API(gin *gin.RouterGroup) {

	a := gin.Group("apikey")
	a.Use(
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		a.POST("create", apiKeyController.Create)
		a.POST("update", apiKeyController.Update)
		a.POST("delete", apiKeyController.Delete)
	}

	b := gin.Group("apikey")
	b.Use(
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		b.GET("list", apiKeyController.List)
		b.GET("get", apiKeyController.Get)
	}
}

func (apiKeyController apiKeyController) Create(ctx *gin.Context) {
	r := new(types.RequestApiKeyCreate)

	// 从JWT中获取当前用户ID
	token := ctx.Request.Header.Get("Authorization")
	userId := tools.GetUserID(token)
	r.UserId = userId

	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.ApiKeyService.Create(r)
	})
}

func (apiKeyController apiKeyController) List(ctx *gin.Context) {
	r := new(types.RequestApiKeyQuery)

	// 从JWT中获取当前用户ID，确保用户只能看到自己的API密钥
	token := ctx.Request.Header.Get("Authorization")
	userId := tools.GetUserID(token)
	r.UserId = userId

	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.ApiKeyService.List(r)
	})
}

func (apiKeyController apiKeyController) Get(ctx *gin.Context) {
	r := new(types.RequestApiKeyQuery)

	// 从JWT中获取当前用户ID，确保用户只能获取自己的API密钥
	token := ctx.Request.Header.Get("Authorization")
	userId := tools.GetUserID(token)
	r.UserId = userId

	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.ApiKeyService.Get(r)
	})
}

func (apiKeyController apiKeyController) Update(ctx *gin.Context) {
	r := new(types.RequestApiKeyUpdate)

	// 从JWT中获取当前用户ID，确保用户只能更新自己的API密钥
	token := ctx.Request.Header.Get("Authorization")
	userId := tools.GetUserID(token)
	r.UserId = userId

	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.ApiKeyService.Update(r)
	})
}

func (apiKeyController apiKeyController) Delete(ctx *gin.Context) {
	r := new(types.RequestApiKeyQuery)

	// 从JWT中获取当前用户ID，确保用户只能删除自己的API密钥
	token := ctx.Request.Header.Get("Authorization")
	userId := tools.GetUserID(token)
	r.UserId = userId

	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.ApiKeyService.Delete(r)
	})
}
