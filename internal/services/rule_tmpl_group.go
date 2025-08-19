package services

import (
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
)

type ruleTmplGroupService struct {
	ctx *ctx.Context
}

type InterRuleTmplGroupService interface {
	List(req interface{}) (interface{}, interface{})
	Create(req interface{}) (interface{}, interface{})
	Update(req interface{}) (interface{}, interface{})
	Delete(req interface{}) (interface{}, interface{})
}

func newInterRuleTmplGroupService(ctx *ctx.Context) InterRuleTmplGroupService {
	return &ruleTmplGroupService{
		ctx: ctx,
	}
}

func (rtg ruleTmplGroupService) List(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRuleTemplateGroupQuery)
	data, count, err := rtg.ctx.DB.RuleTmplGroup().List(r.Type, r.Query, r.Page)
	if err != nil {
		return nil, err
	}

	return types.ResponseRuleTemplateGroupList{
		List: data,
		Page: models.Page{
			Total: count,
			Index: r.Page.Index,
			Size:  r.Page.Size,
		},
	}, nil
}

func (rtg ruleTmplGroupService) Create(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRuleTemplateGroupCreate)
	err := rtg.ctx.DB.RuleTmplGroup().Create(models.RuleTemplateGroup{
		Name:        r.Name,
		Type:        r.Type,
		Description: r.Description,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (rtg ruleTmplGroupService) Update(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRuleTemplateGroupUpdate)
	err := rtg.ctx.DB.RuleTmplGroup().Update(models.RuleTemplateGroup{
		Name:        r.Name,
		Type:        r.Type,
		Description: r.Description,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (rtg ruleTmplGroupService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRuleTemplateGroupQuery)
	err := rtg.ctx.DB.RuleTmplGroup().Delete(r.Name)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
