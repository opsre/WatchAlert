package services

import (
	"fmt"
	"watchAlert/internal/ctx"
	"watchAlert/internal/types"
	"watchAlert/pkg/provider"
)

type (
	clientService struct {
		ctx *ctx.Context
	}

	InterClientService interface {
		GetJaegerService(req interface{}) (interface{}, interface{})
		GetTraceService(req interface{}) (interface{}, interface{})
	}
)

func newInterClientService(ctx *ctx.Context) InterClientService {
	return &clientService{
		ctx: ctx,
	}
}

func (cs clientService) GetJaegerService(req interface{}) (interface{}, interface{}) {
	return cs.GetTraceService(req)
}

func (cs clientService) GetTraceService(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestDatasourceQuery)

	getInfo, err := cs.ctx.DB.Datasource().Get(r.TenantId, r.ID)
	if err != nil {
		return nil, err
	}
	if getInfo.Type != provider.JaegerDsProviderName && getInfo.Type != provider.SkyWalkingDsProviderName {
		return nil, fmt.Errorf("数据源类型 %s 不支持链路服务查询", getInfo.Type)
	}

	cli, err := provider.NewTraceClient(getInfo)
	if err != nil {
		return nil, err
	}
	service, err := cli.GetServiceNames()
	if err != nil {
		return nil, err
	}

	return service, nil
}
