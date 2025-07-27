package services

import (
	"watchAlert/internal/ctx"
	"watchAlert/internal/types"
)

type auditLogService struct {
	ctx *ctx.Context
}

type InterAuditLogService interface {
	List(req interface{}) (interface{}, interface{})
	Search(req interface{}) (interface{}, interface{})
}

func newInterAuditLogService(ctx *ctx.Context) InterAuditLogService {
	return &auditLogService{
		ctx: ctx,
	}
}

func (as auditLogService) List(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestAuditLogQuery)
	data, err := as.ctx.DB.AuditLog().List(*r)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (as auditLogService) Search(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestAuditLogQuery)
	data, err := as.ctx.DB.AuditLog().Search(*r)
	if err != nil {
		return nil, err
	}

	return data, nil
}
