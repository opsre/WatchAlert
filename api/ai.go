package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
)

type AiController struct{}

func (a AiController) API(gin *gin.RouterGroup) {
	ai := gin.Group("ai")
	ai.Use(
		middleware.Cors(),
		middleware.Auth(),
	)
	{
		ai.POST("chat", a.Chat)
	}
}

func (a AiController) Chat(ctx *gin.Context) {
	r := new(types.RequestAiChatContent)
	r.Content = ctx.PostForm("content")
	r.RuleId = ctx.PostForm("rule_id")
	r.RuleName = ctx.PostForm("rule_name")
	r.Deep = ctx.PostForm("deep")
	r.SearchQL = ctx.PostForm("search_ql")

	Service(ctx, func() (interface{}, interface{}) {
		return services.AiService.Chat(r)
	})
}
