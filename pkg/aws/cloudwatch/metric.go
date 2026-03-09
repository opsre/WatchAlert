package cloudwatch

import (
	"watchAlert/internal/ctx"
)

type (
	awsCloudWatchService struct {
		ctx *ctx.Context
	}

	InterAwsCloudWatchService interface {
		GetMetricTypes() (interface{}, interface{})
		GetMetricNames(req interface{}) (interface{}, interface{})
		GetStatistics() (interface{}, interface{})
		GetDimensions(req interface{}) (interface{}, interface{})
	}
)

func NewInterAwsCloudWatchService(ctx *ctx.Context) InterAwsCloudWatchService {
	return awsCloudWatchService{
		ctx: ctx,
	}
}

func (a awsCloudWatchService) GetMetricTypes() (interface{}, interface{}) {
	var mt []string
	for k, _ := range NamespaceMetricsMap {
		mt = append(mt, k)
	}

	return mt, nil
}

func (a awsCloudWatchService) GetMetricNames(req interface{}) (interface{}, interface{}) {
	r := req.(*MetricNamesQuery)
	return NamespaceMetricsMap[r.MetricType], nil
}

func (a awsCloudWatchService) GetStatistics() (interface{}, interface{}) {
	return []string{
		"Average",
		"Maximum",
		"Minimum",
		"Sum",
		"SampleCount",
		"IQM",
	}, nil
}

func (a awsCloudWatchService) GetDimensions(req interface{}) (interface{}, interface{}) {
	r := req.(*RdsDimensionReq)
	return NamespaceDimensionKeysMap[r.MetricType], nil
}
