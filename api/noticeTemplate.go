package api

import (
	"errors"
	"watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
	"watchAlert/pkg/tools"

	"github.com/gin-gonic/gin"
)

type noticeTemplateController struct{}

var NoticeTemplateController = new(noticeTemplateController)

/*
通知模版 API
/api/w8t/noticeTemplate
*/
func (noticeTemplateController noticeTemplateController) API(gin *gin.RouterGroup) {
	a := gin.Group("noticeTemplate")
	a.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		a.POST("noticeTemplateCreate", noticeTemplateController.Create)
		a.POST("noticeTemplateUpdate", noticeTemplateController.Update)
		a.POST("noticeTemplateDelete", noticeTemplateController.Delete)
	}
	b := gin.Group("noticeTemplate")
	b.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		b.GET("noticeTemplateList", noticeTemplateController.List)
	}
}

func (noticeTemplateController noticeTemplateController) Create(ctx *gin.Context) {
	r := new(types.RequestNoticeTemplateCreate)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		tokenStr := ctx.Request.Header.Get("Authorization")
		if len(tokenStr) <= 0 {
			return nil, errors.New("用户未登录")
		}
		r.UpdateBy = tools.GetUser(tokenStr)

		return services.NoticeTmplService.Create(r)
	})
}

func (noticeTemplateController noticeTemplateController) Update(ctx *gin.Context) {
	r := new(types.RequestNoticeTemplateUpdate)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		tokenStr := ctx.Request.Header.Get("Authorization")
		if len(tokenStr) <= 0 {
			return nil, errors.New("用户未登录")
		}
		r.UpdateBy = tools.GetUser(tokenStr)

		return services.NoticeTmplService.Update(r)
	})
}

func (noticeTemplateController noticeTemplateController) Delete(ctx *gin.Context) {
	r := new(types.RequestNoticeTemplateQuery)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.NoticeTmplService.Delete(r)
	})
}

func (noticeTemplateController noticeTemplateController) List(ctx *gin.Context) {
	r := new(types.RequestNoticeTemplateQuery)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.NoticeTmplService.List(r)
	})
}
