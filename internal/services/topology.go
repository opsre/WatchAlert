package services

import (
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
	"watchAlert/pkg/tools"
)

type (
	topologyService struct {
		ctx *ctx.Context
	}

	InterTopologyService interface {
		Create(req interface{}) (data interface{}, err interface{})
		Update(req interface{}) (data interface{}, err interface{})
		Delete(req interface{}) (data interface{}, err interface{})
		List(req interface{}) (data interface{}, err interface{})
		GetDetail(req interface{}) (data interface{}, err interface{})
	}
)

func newInterTopologyService(ctx *ctx.Context) InterTopologyService {
	return &topologyService{
		ctx: ctx,
	}
}

func (t topologyService) Create(req interface{}) (data interface{}, err interface{}) {
	r := req.(*types.RequestTopologyCreate)
	topology := models.Topology{
		TenantId:  r.TenantId,
		ID:        "topo-" + tools.RandId(),
		Name:      r.Name,
		Nodes:     r.Nodes,
		Edges:     r.Edges,
		UpdatedBy: r.UpdatedBy,
		UpdatedAt: time.Now().Unix(),
	}

	err = t.ctx.DB.Topology().Create(topology)
	if err != nil {
		return nil, err
	}

	t.ctx.Redis.Topology().PushTopologyInfo(topology)

	return nil, nil
}

func (t topologyService) Update(req interface{}) (data interface{}, err interface{}) {
	r := req.(*types.RequestTopologyUpdate)
	topology := models.Topology{
		TenantId:  r.TenantId,
		ID:        r.ID,
		Name:      r.Name,
		Nodes:     r.Nodes,
		Edges:     r.Edges,
		UpdatedBy: r.UpdatedBy,
		UpdatedAt: time.Now().Unix(),
	}

	err = t.ctx.DB.Topology().Update(topology)
	if err != nil {
		return nil, err
	}

	t.ctx.Redis.Topology().PushTopologyInfo(topology)

	return nil, nil
}

func (t topologyService) Delete(req interface{}) (data interface{}, err interface{}) {
	r := req.(*types.RequestTopologyDelete)
	err = t.ctx.DB.Topology().Delete(r.TenantId, r.ID)
	if err != nil {
		return nil, err
	}

	t.ctx.Redis.Topology().RemoveTopologyInfo(models.BuildTopologyCacheKey(r.TenantId, r.ID))

	return nil, nil
}

func (t topologyService) List(req interface{}) (data interface{}, err interface{}) {
	r := req.(*types.RequestTopologyQuery)
	data, err = t.ctx.DB.Topology().List(r.TenantId, r.Query)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetDetail 获取拓扑的完整信息，包括nodes和edges
func (t topologyService) GetDetail(req interface{}) (data interface{}, err interface{}) {
	r := req.(*types.RequestTopologyQuery)
	data, err = t.ctx.DB.Topology().GetDetail(r.TenantId, r.ID)
	if err != nil {
		return nil, err
	}
	return data, nil
}
