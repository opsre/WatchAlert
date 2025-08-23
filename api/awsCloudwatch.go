package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/pkg/community/aws/cloudwatch/types"
)

type awsCloudWatchController struct{}

var AWSCloudWatchController = new(awsCloudWatchController)

func (awsCloudWatchController awsCloudWatchController) API(gin *gin.RouterGroup) {
	community := gin.Group("community")
	community.Use(
		middleware.Cors(),
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		cloudwatch := community.Group("cloudwatch")
		{
			cloudwatch.GET("metricTypes", awsCloudWatchController.GetMetricTypes)
			cloudwatch.GET("metricNames", awsCloudWatchController.GetMetricNames)
			cloudwatch.GET("statistics", awsCloudWatchController.GetStatistics)
			cloudwatch.GET("dimensions", awsCloudWatchController.GetDimensions)
		}
	}
}

func (awsCloudWatchController awsCloudWatchController) GetMetricTypes(ctx *gin.Context) {
	Service(ctx, func() (interface{}, interface{}) {
		return services.AWSCloudWatchService.GetMetricTypes()
	})
}

func (awsCloudWatchController awsCloudWatchController) GetMetricNames(ctx *gin.Context) {
	q := new(types.MetricNamesQuery)
	BindQuery(ctx, q)
	Service(ctx, func() (interface{}, interface{}) {
		return services.AWSCloudWatchService.GetMetricNames(q)
	})
}

func (awsCloudWatchController awsCloudWatchController) GetStatistics(ctx *gin.Context) {
	Service(ctx, func() (interface{}, interface{}) {
		return services.AWSCloudWatchService.GetStatistics()
	})
}

func (awsCloudWatchController awsCloudWatchController) GetDimensions(ctx *gin.Context) {
	q := new(types.RdsDimensionReq)
	BindQuery(ctx, q)
	Service(ctx, func() (interface{}, interface{}) {
		return services.AWSCloudWatchService.GetDimensions(q)
	})
}
