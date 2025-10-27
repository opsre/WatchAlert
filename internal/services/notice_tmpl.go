package services

import (
	"errors"
	"fmt"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
	"watchAlert/pkg/tools"
)

type noticeTmplService struct {
	ctx *ctx.Context
}

type InterNoticeTmplService interface {
	List(req interface{}) (interface{}, interface{})
	Create(req interface{}) (interface{}, interface{})
	Update(req interface{}) (interface{}, interface{})
	Delete(req interface{}) (interface{}, interface{})
}

func newInterNoticeTmplService(ctx *ctx.Context) InterNoticeTmplService {
	return &noticeTmplService{
		ctx,
	}
}

func (nts noticeTmplService) List(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestNoticeTemplateQuery)
	data, err := nts.ctx.DB.NoticeTmpl().List(r.ID, r.NoticeType, r.Query)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (nts noticeTmplService) Create(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestNoticeTemplateCreate)
	err := nts.ctx.DB.NoticeTmpl().Create(models.NoticeTemplateExample{
		ID:                   "nt-" + tools.RandId(),
		Name:                 r.Name,
		NoticeType:           r.NoticeType,
		Description:          r.Description,
		Template:             r.Template,
		TemplateFiring:       r.TemplateFiring,
		TemplateRecover:      r.TemplateRecover,
		EnableFeiShuJsonCard: r.EnableFeiShuJsonCard,
		UpdateAt:             time.Now().Unix(),
		UpdateBy:             r.UpdateBy,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (nts noticeTmplService) Update(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestNoticeTemplateUpdate)
	err := nts.ctx.DB.NoticeTmpl().Update(models.NoticeTemplateExample{
		ID:                   r.ID,
		Name:                 r.Name,
		NoticeType:           r.NoticeType,
		Description:          r.Description,
		Template:             r.Template,
		TemplateFiring:       r.TemplateFiring,
		TemplateRecover:      r.TemplateRecover,
		EnableFeiShuJsonCard: r.EnableFeiShuJsonCard,
		UpdateAt:             time.Now().Unix(),
		UpdateBy:             r.UpdateBy,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (nts noticeTmplService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestNoticeTemplateQuery)
	nl, err := nts.ctx.DB.Notice().List("", r.ID, "")
	if err != nil {
		return nil, err
	}

	if len(nl) > 0 {
		var ids []string
		for _, n := range nl {
			ids = append(ids, n.Uuid)
		}
		return nil, errors.New(fmt.Sprintf("删除失败, 已有通知对象绑定: %s", ids))
	}

	err = nts.ctx.DB.NoticeTmpl().Delete(r.ID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
