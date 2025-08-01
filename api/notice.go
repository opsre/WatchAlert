package api

import (
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
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
		b.GET("noticeRecordMetric", noticeController.GetRecordMetric)
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

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.NoticeService.Create(r)
	})
}

func (noticeController noticeController) Update(ctx *gin.Context) {
	r := new(types.RequestNoticeUpdate)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
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

func (noticeController noticeController) Check(ctx *gin.Context) {
	r := new(types.RequestNoticeQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.NoticeService.Check(r)
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
