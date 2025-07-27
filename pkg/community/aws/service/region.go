package service

import (
	"watchAlert/internal/ctx"
	"watchAlert/pkg/community/aws/cloudwatch/response"
	"watchAlert/pkg/community/aws/types"
)

type (
	awsRegionService struct {
		ctx *ctx.Context
	}

	InterAwsRegionService interface {
		GetRegion() ([]response.Regions, error)
	}
)

func NewInterAwsRegionService(ctx *ctx.Context) InterAwsRegionService {
	return awsRegionService{
		ctx: ctx,
	}
}

func (a awsRegionService) GetRegion() ([]response.Regions, error) {
	var rs []response.Regions
	for _, r := range types.Regions {
		rs = append(rs, response.Regions{
			Label: &r,
			Value: &r,
		})
	}

	return rs, nil
}
