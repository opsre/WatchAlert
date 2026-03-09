package services

import (
	"errors"
	"fmt"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
	"watchAlert/pkg/sender"
	"watchAlert/pkg/tools"

	"github.com/zeromicro/go-zero/core/logc"
)

type noticeService struct {
	ctx *ctx.Context
}

type InterNoticeService interface {
	List(req interface{}) (interface{}, interface{})
	Create(req interface{}) (interface{}, interface{})
	Update(req interface{}) (interface{}, interface{})
	Delete(req interface{}) (interface{}, interface{})
	Get(req interface{}) (interface{}, interface{})
	ListRecord(req interface{}) (interface{}, interface{})
	GetRecordMetric(req interface{}) (interface{}, interface{})
	DeleteRecord(req interface{}) (interface{}, interface{})
	Test(req interface{}) (interface{}, interface{})
}

func newInterAlertNoticeService(ctx *ctx.Context) InterNoticeService {
	return &noticeService{
		ctx,
	}
}

func (n noticeService) List(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestNoticeQuery)
	data, err := n.ctx.DB.Notice().List(r.TenantId, r.NoticeTmplId, r.Query)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (n noticeService) Create(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestNoticeCreate)
	ok := n.ctx.DB.Notice().GetQuota(r.TenantId)
	if !ok {
		return models.AlertNotice{}, fmt.Errorf("创建失败, 配额不足")
	}

	err := n.ctx.DB.Notice().Create(models.AlertNotice{
		TenantId: r.TenantId,
		Uuid:     "n-" + tools.RandId(),
		Name:     r.Name,
		DutyId:   r.DutyId,
		Routes:   r.Routes,
		UpdateAt: time.Now().Unix(),
		UpdateBy: r.UpdateBy,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (n noticeService) Update(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestNoticeUpdate)
	err := n.ctx.DB.Notice().Update(models.AlertNotice{
		TenantId: r.TenantId,
		Uuid:     r.Uuid,
		Name:     r.Name,
		DutyId:   r.GetDutyId(),
		Routes:   r.Routes,
		UpdateAt: time.Now().Unix(),
		UpdateBy: r.UpdateBy,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (n noticeService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestNoticeQuery)
	err := n.ctx.DB.Notice().Delete(r.TenantId, r.Uuid)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (n noticeService) Get(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestNoticeQuery)
	data, err := n.ctx.DB.Notice().Get(r.TenantId, r.Uuid)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (n noticeService) ListRecord(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestNoticeQuery)
	data, err := n.ctx.DB.Notice().ListRecord(r.TenantId, r.EventId, r.Severity, r.Status, r.Uuid, r.Query, r.Page)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (n noticeService) DeleteRecord(req interface{}) (interface{}, interface{}) {
	err := n.ctx.DB.Notice().DeleteRecord()
	if err != nil {
		return nil, err
	}

	return nil, nil
}

type ResponseRecordMetric struct {
	Date   []string `json:"date"`
	Series series   `json:"series"`
}

type series struct {
	P0 []int64 `json:"p0"`
	P1 []int64 `json:"p1"`
	P2 []int64 `json:"p2"`
}

func (n noticeService) GetRecordMetric(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestNoticeQuery)
	curTime := time.Now()
	var layout = "2006-01-02"
	timeList := []string{
		curTime.Add(-144 * time.Hour).Format(layout),
		curTime.Add(-120 * time.Hour).Format(layout),
		curTime.Add(-96 * time.Hour).Format(layout),
		curTime.Add(-72 * time.Hour).Format(layout),
		curTime.Add(-48 * time.Hour).Format(layout),
		curTime.Add(-24 * time.Hour).Format(layout),
		curTime.Format(layout),
	}

	var severitys = []string{"P0", "P1", "P2"}
	var P0, P1, P2 []int64
	for _, t := range timeList {
		for _, s := range severitys {
			count, err := n.ctx.DB.Notice().CountRecord(models.CountRecord{
				Date:     t,
				TenantId: r.TenantId,
				Severity: s,
			})
			if err != nil {
				logc.Error(n.ctx.Ctx, err.Error())
			}
			switch s {
			case "P0":
				P0 = append(P0, count)
			case "P1":
				P1 = append(P1, count)
			case "P2":
				P2 = append(P2, count)
			}

		}
	}

	return ResponseRecordMetric{
		Date: timeList,
		Series: series{
			P0: P0,
			P1: P1,
			P2: P2,
		},
	}, nil
}

func (n noticeService) Test(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestNoticeTest)
	var errList []struct {
		Hook  string
		Error string
	}

	err := sender.Tester(n.ctx, sender.SendParams{
		NoticeType: r.NoticeType,
		Hook:       r.Hook,
		Email:      r.Email,
		Sign:       r.Sign,
	})
	if err != nil {
		errList = append(errList, struct {
			Hook  string
			Error string
		}{Hook: r.Hook, Error: err.Error()})
	}

	if len(errList) != 0 {
		return nil, errors.New(tools.JsonMarshalToString(errList))
	}

	return nil, nil
}
