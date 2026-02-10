package region

import (
	"watchAlert/internal/ctx"
)

type (
	awsRegionService struct {
		ctx *ctx.Context
	}

	InterAwsRegionService interface {
		GetRegion() ([]RegionItem, error)
	}
)

func NewInterAwsRegionService(ctx *ctx.Context) InterAwsRegionService {
	return awsRegionService{
		ctx: ctx,
	}
}

func (a awsRegionService) GetRegion() ([]RegionItem, error) {
	var rs []RegionItem
	for _, r := range Regions {
		rs = append(rs, RegionItem{
			Label: &r,
			Value: &r,
		})
	}

	return rs, nil
}
