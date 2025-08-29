package services

import (
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
)

type ruleTmplService struct {
	ctx *ctx.Context
}

type InterRuleTmplService interface {
	List(req interface{}) (interface{}, interface{})
	Create(req interface{}) (interface{}, interface{})
	Update(req interface{}) (interface{}, interface{})
	Delete(req interface{}) (interface{}, interface{})
}

func newInterRuleTmplService(ctx *ctx.Context) InterRuleTmplService {
	return &ruleTmplService{
		ctx: ctx,
	}
}

func (rt ruleTmplService) List(req interface{}) (interface{}, interface{}) {
	r := req.(*models.RuleTemplateQuery)
	data, err := rt.ctx.DB.RuleTmpl().List(*r)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (rt ruleTmplService) Create(req interface{}) (interface{}, interface{}) {
	r := req.(*models.RuleTemplate)
	err := rt.ctx.DB.RuleTmpl().Create(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (rt ruleTmplService) Update(req interface{}) (interface{}, interface{}) {
	r := req.(*models.RuleTemplate)
	err := rt.ctx.DB.RuleTmpl().Update(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (rt ruleTmplService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*models.RuleTemplateQuery)
	err := rt.ctx.DB.RuleTmpl().Delete(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
