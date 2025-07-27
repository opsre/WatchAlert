package services

import (
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/pkg/provider"
)

type (
	clientService struct {
		ctx *ctx.Context
	}

	InterClientService interface {
		GetJaegerService(req interface{}) (interface{}, interface{})
	}
)

func newInterClientService(ctx *ctx.Context) InterClientService {
	return &clientService{
		ctx: ctx,
	}
}

func (cs clientService) GetJaegerService(req interface{}) (interface{}, interface{}) {
	r := req.(*models.DatasourceQuery)

	getInfo, err := cs.ctx.DB.Datasource().Get(*r)
	if err != nil {
		return nil, err
	}

	cli, err := provider.NewJaegerClient(getInfo)
	service, err := cli.GetJaegerService()
	if err != nil {
		return nil, err
	}

	return service, nil
}
