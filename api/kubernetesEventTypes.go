package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/internal/middleware"
	"watchAlert/internal/types"
)

type kubernetesTypesController struct{}

var KubernetesTypesController = new(kubernetesTypesController)

func (kubernetesTypesController kubernetesTypesController) API(gin *gin.RouterGroup) {
	k8s := gin.Group("kubernetes")
	k8s.Use(
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		k8s.GET("getResourceList", kubernetesTypesController.getResourceList)
		k8s.GET("getReasonList", kubernetesTypesController.getReasonList)
	}
}

func (kubernetesTypesController kubernetesTypesController) getResourceList(ctx *gin.Context) {
	Service(ctx, func() (interface{}, interface{}) {
		return types.EventResourceTypeList, nil
	})
}

func (kubernetesTypesController kubernetesTypesController) getReasonList(ctx *gin.Context) {
	r := new(types.RequestKubernetesEventTypes)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return types.EventReasonLMapping[r.Resource], nil
	})
}
