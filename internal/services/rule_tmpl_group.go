package services

import (
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
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
	r := req.(*models.RuleTemplateGroupQuery)
	data, err := rtg.ctx.DB.RuleTmplGroup().List(*r)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (rtg ruleTmplGroupService) Create(req interface{}) (interface{}, interface{}) {
	r := req.(*models.RuleTemplateGroup)
	err := rtg.ctx.DB.RuleTmplGroup().Create(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (rtg ruleTmplGroupService) Update(req interface{}) (interface{}, interface{}) {
	r := req.(*models.RuleTemplateGroup)
	err := rtg.ctx.DB.RuleTmplGroup().Update(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (rtg ruleTmplGroupService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*models.RuleTemplateGroupQuery)
	err := rtg.ctx.DB.RuleTmplGroup().Delete(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
