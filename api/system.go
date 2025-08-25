package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
)

type systemController struct{}

var SystemController = new(systemController)

func (systemController systemController) GetOidcInfo(ctx *gin.Context) {
	Service(ctx, func() (interface{}, interface{}) { return services.OidcService.GetOidcInfo() })
}

func (systemController systemController) CallBack(ctx *gin.Context) {
	r := new(types.RequestOidcCodeQuery)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) { return services.OidcService.CallBack(ctx, r) })
}

func (systemController systemController) CookieConvertToken(ctx *gin.Context) {
	Service(ctx, func() (interface{}, interface{}) {
		return services.OidcService.CookieConvertToken(ctx)
	})
}
