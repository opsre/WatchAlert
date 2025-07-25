package services

import (
	"errors"
	"fmt"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
)

type noticeTmplService struct {
	ctx *ctx.Context
}

type InterNoticeTmplService interface {
	List(req interface{}) (interface{}, interface{})
	Search(req interface{}) (interface{}, interface{})
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
	r := req.(*models.NoticeTemplateExampleQuery)
	data, err := nts.ctx.DB.NoticeTmpl().List(*r)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (nts noticeTmplService) Search(req interface{}) (interface{}, interface{}) {
	r := req.(*models.NoticeTemplateExampleQuery)
	data, err := nts.ctx.DB.NoticeTmpl().Search(*r)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (nts noticeTmplService) Create(req interface{}) (interface{}, interface{}) {
	r := req.(*models.NoticeTemplateExample)
	err := nts.ctx.DB.NoticeTmpl().Create(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (nts noticeTmplService) Update(req interface{}) (interface{}, interface{}) {
	r := req.(*models.NoticeTemplateExample)
	err := nts.ctx.DB.NoticeTmpl().Update(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (nts noticeTmplService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*models.NoticeTemplateExampleQuery)
	nl, err := nts.ctx.DB.Notice().Search(models.NoticeQuery{NoticeTmplId: r.Id})
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

	err = nts.ctx.DB.NoticeTmpl().Delete(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
