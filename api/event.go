package api

import (
	"time"
	"watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
	"watchAlert/pkg/response"
	utils "watchAlert/pkg/tools"

	"github.com/gin-gonic/gin"
)

type alertEventController struct{}

var AlertEventController = new(alertEventController)

/*
告警事件 API
/api/w8t/event
*/
func (alertEventController alertEventController) API(gin *gin.RouterGroup) {
	a := gin.Group("event")
	a.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		a.POST("process", alertEventController.ProcessAlertEvent)
		a.POST("delete", alertEventController.DeleteAlertEvent)
		a.POST("addComment", alertEventController.AddComment)
		a.GET("listComments", alertEventController.ListComment)
		a.POST("deleteComment", alertEventController.DeleteComment)
	}

	b := gin.Group("event")
	b.Use(
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		b.GET("curEvent", alertEventController.ListCurrentEvent)
		b.GET("hisEvent", alertEventController.ListHistoryEvent)
	}
}

func (alertEventController alertEventController) ProcessAlertEvent(ctx *gin.Context) {
	r := new(types.RequestProcessAlertEvent)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)
	r.Time = time.Now().Unix()

	tokenStr := ctx.Request.Header.Get("Authorization")
	if tokenStr == "" {
		response.Fail(ctx, "未知的用户", "failed")
		return
	}

	r.Username = utils.GetUser(tokenStr)

	Service(ctx, func() (interface{}, interface{}) {
		return services.EventService.ProcessAlertEvent(r)
	})
}

func (alertEventController alertEventController) DeleteAlertEvent(ctx *gin.Context) {
	r := new(types.RequestProcessAlertEvent)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)
	r.Time = time.Now().Unix()

	Service(ctx, func() (interface{}, interface{}) {
		return services.EventService.DeleteAlertEvent(r)
	})
}

func (alertEventController alertEventController) ListCurrentEvent(ctx *gin.Context) {
	r := new(types.RequestAlertCurEventQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.EventService.ListCurrentEvent(r)
	})
}

func (alertEventController alertEventController) ListHistoryEvent(ctx *gin.Context) {
	r := new(types.RequestAlertHisEventQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.EventService.ListHistoryEvent(r)
	})
}

func (alertEventController alertEventController) ListComment(ctx *gin.Context) {
	r := new(types.RequestListEventComments)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.EventService.ListComments(r)
	})
}

func (alertEventController alertEventController) AddComment(ctx *gin.Context) {
	r := new(types.RequestAddEventComment)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	token := ctx.Request.Header.Get("Authorization")
	r.Username = utils.GetUser(token)
	r.UserId = utils.GetUserID(token)

	Service(ctx, func() (interface{}, interface{}) {
		return services.EventService.AddComment(r)
	})
}

func (alertEventController alertEventController) DeleteComment(ctx *gin.Context) {
	r := new(types.RequestDeleteEventComment)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.EventService.DeleteComment(r)
	})
}
