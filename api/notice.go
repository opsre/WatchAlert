package api

import (
	"errors"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
	"watchAlert/pkg/tools"

	"github.com/gin-gonic/gin"
)

type noticeController struct{}

var NoticeController = new(noticeController)

/*
通知对象 API
/api/w8t/sender
*/
func (noticeController noticeController) API(gin *gin.RouterGroup) {
	a := gin.Group("notice")
	a.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		a.POST("noticeCreate", noticeController.Create)
		a.POST("noticeUpdate", noticeController.Update)
		a.POST("noticeDelete", noticeController.Delete)
	}

	b := gin.Group("notice")
	b.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		b.GET("noticeList", noticeController.List)
		b.GET("noticeRecordList", noticeController.ListRecord)
	}

	c := gin.Group("notice")
	c.Use(
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		c.GET("noticeRecordMetric", noticeController.GetRecordMetric)
		c.POST("noticeTest", noticeController.Test)
	}
}

func (noticeController noticeController) List(ctx *gin.Context) {
	r := new(types.RequestNoticeQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.NoticeService.List(r)
	})
}

func (noticeController noticeController) Create(ctx *gin.Context) {
	r := new(types.RequestNoticeCreate)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		tokenStr := ctx.Request.Header.Get("Authorization")
		if len(tokenStr) <= 0 {
			return nil, errors.New("用户未登录")
		}
		r.UpdateBy = tools.GetUser(tokenStr)

		tid, _ := ctx.Get("TenantID")
		r.TenantId = tid.(string)

		return services.NoticeService.Create(r)
	})
}

func (noticeController noticeController) Update(ctx *gin.Context) {
	r := new(types.RequestNoticeUpdate)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		tokenStr := ctx.Request.Header.Get("Authorization")
		if len(tokenStr) <= 0 {
			return nil, errors.New("用户未登录")
		}
		r.UpdateBy = tools.GetUser(tokenStr)

		tid, _ := ctx.Get("TenantID")
		r.TenantId = tid.(string)

		return services.NoticeService.Update(r)
	})
}

func (noticeController noticeController) Delete(ctx *gin.Context) {
	r := new(types.RequestNoticeQuery)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.NoticeService.Delete(r)
	})
}

func (noticeController noticeController) Get(ctx *gin.Context) {
	r := new(types.RequestNoticeQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.NoticeService.Get(r)
	})

}

func (noticeController noticeController) ListRecord(ctx *gin.Context) {
	r := new(types.RequestNoticeQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.NoticeService.ListRecord(r)
	})
}

func (noticeController noticeController) GetRecordMetric(ctx *gin.Context) {
	r := new(types.RequestNoticeQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.NoticeService.GetRecordMetric(r)
	})
}

func (noticeController noticeController) Test(ctx *gin.Context) {
	r := new(types.RequestNoticeTest)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.NoticeService.Test(r)
	})
}
