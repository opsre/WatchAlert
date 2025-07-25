package services

import (
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
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
	r := req.(*models.RuleGroups)
	err := rgs.ctx.DB.RuleGroup().Create(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (rgs ruleGroupService) Update(req interface{}) (interface{}, interface{}) {
	r := req.(*models.RuleGroups)
	err := rgs.ctx.DB.RuleGroup().Update(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (rgs ruleGroupService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*models.RuleGroupQuery)
	err := rgs.ctx.DB.RuleGroup().Delete(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (rgs ruleGroupService) List(req interface{}) (interface{}, interface{}) {
	r := req.(*models.RuleGroupQuery)
	data, err := rgs.ctx.DB.RuleGroup().List(*r)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (rgs ruleGroupService) Search() {

}
