package services

import (
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
	"watchAlert/pkg/tools"
)

type ruleGroupService struct {
	ctx *ctx.Context
}

type InterRuleGroupService interface {
	Create(req interface{}) (interface{}, interface{})
	Update(req interface{}) (interface{}, interface{})
	Delete(req interface{}) (interface{}, interface{})
	List(req interface{}) (interface{}, interface{})
}

func newInterRuleGroupService(ctx *ctx.Context) InterRuleGroupService {
	return &ruleGroupService{
		ctx: ctx,
	}
}

func (rgs ruleGroupService) Create(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRuleGroupCreate)
	err := rgs.ctx.DB.RuleGroup().Create(models.RuleGroups{
		TenantId:    r.TenantId,
		ID:          "rg-" + tools.RandId(),
		Name:        r.Name,
		Description: r.Description,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (rgs ruleGroupService) Update(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRuleGroupUpdate)
	err := rgs.ctx.DB.RuleGroup().Update(models.RuleGroups{
		TenantId:    r.TenantId,
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (rgs ruleGroupService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRuleGroupQuery)
	err := rgs.ctx.DB.RuleGroup().Delete(r.TenantId, r.ID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (rgs ruleGroupService) List(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRuleGroupQuery)
	data, count, err := rgs.ctx.DB.RuleGroup().List(r.TenantId, r.Query, r.Page)
	if err != nil {
		return nil, err
	}

	return types.ResponseRuleGroupList{
		List: data,
		Page: models.Page{
			Index: r.Page.Index,
			Size:  r.Page.Size,
			Total: count,
		},
	}, nil
}
