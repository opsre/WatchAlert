package services

import (
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
)

type recordingRuleGroupService struct {
	ctx *ctx.Context
}

type InterRecordingRuleGroupService interface {
	Create(req interface{}) (interface{}, interface{})
	Update(req interface{}) (interface{}, interface{})
	Delete(req interface{}) (interface{}, interface{})
	List(req interface{}) (interface{}, interface{})
	Get(req interface{}) (interface{}, interface{})
}

func newInterRecordingRuleGroupService(ctx *ctx.Context) InterRecordingRuleGroupService {
	return &recordingRuleGroupService{
		ctx: ctx,
	}
}

func (rgs recordingRuleGroupService) Create(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRecordingRuleGroupCreate)
	err := rgs.ctx.DB.RecordingRuleGroup().Create(&models.RecordingRuleGroup{
		TenantId: r.TenantId,
		Name:     r.Name,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (rgs recordingRuleGroupService) Update(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRecordingRuleGroupUpdate)
	err := rgs.ctx.DB.RecordingRuleGroup().Update(&models.RecordingRuleGroup{
		ID:       r.ID,
		TenantId: r.TenantId,
		Name:     r.Name,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (rgs recordingRuleGroupService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRecordingRuleGroupQuery)
	err := rgs.ctx.DB.RecordingRuleGroup().Delete(r.TenantId, r.ID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (rgs recordingRuleGroupService) List(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRecordingRuleGroupQuery)
	data, count, err := rgs.ctx.DB.RecordingRuleGroup().List(r.TenantId, r.Query, r.Page)
	if err != nil {
		return nil, err
	}

	return types.ResponseRecordingRuleGroupList{
		List: data,
		Page: models.Page{
			Index: r.Page.Index,
			Size:  r.Page.Size,
			Total: count,
		},
	}, nil
}

func (rgs recordingRuleGroupService) Get(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRecordingRuleGroupQuery)
	data, err := rgs.ctx.DB.RecordingRuleGroup().Get(r.TenantId, r.ID)
	if err != nil {
		return nil, err
	}

	return data, nil
}
