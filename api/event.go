package api

import (
	"github.com/gin-gonic/gin"
	"time"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/services"
	"watchAlert/pkg/response"
	utils "watchAlert/pkg/tools"
)

type AlertEventController struct{}

/*
告警事件 API
/api/w8t/event
*/
func (e AlertEventController) API(gin *gin.RouterGroup) {
	event := gin.Group("event")
	event.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		event.GET("curEvent", e.ListCurrentEvent)
		event.GET("hisEvent", e.ListHistoryEvent)
		event.POST("processAlertEvent", e.ProcessAlertEvent)
		event.POST("addComment", e.AddComment)
		event.GET("listComments", e.ListComment)
		event.POST("deleteComment", e.DeleteComment)
	}
}

func (e AlertEventController) ProcessAlertEvent(ctx *gin.Context) {
	r := new(models.ProcessAlertEvent)
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

func (e AlertEventController) ListCurrentEvent(ctx *gin.Context) {
	r := new(models.AlertCurEventQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.EventService.ListCurrentEvent(r)
	})
}

func (e AlertEventController) ListHistoryEvent(ctx *gin.Context) {
	r := new(models.AlertHisEventQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.EventService.ListHistoryEvent(r)
	})
}

func (e AlertEventController) ListComment(ctx *gin.Context) {
	r := new(models.RequestListEventComments)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.EventService.ListComments(r)
	})
}

func (e AlertEventController) AddComment(ctx *gin.Context) {
	r := new(models.RequestAddEventComment)
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

func (e AlertEventController) DeleteComment(ctx *gin.Context) {
	r := new(models.RequestDeleteEventComment)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.EventService.DeleteComment(r)
	})
}
