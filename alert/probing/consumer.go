package probing

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"time"
	"watchAlert/alert/process"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/sender"
	"watchAlert/pkg/templates"
	"watchAlert/pkg/tools"
)

type ConsumeProbing struct {
	ctx          *ctx.Context
	consumerPool map[string]context.CancelFunc
}

func NewProbingConsumerTask(ctx *ctx.Context) ConsumeProbing {
	return ConsumeProbing{
		ctx:          ctx,
		consumerPool: make(map[string]context.CancelFunc),
	}
}

func (m *ConsumeProbing) Add(r models.ProbingRule) {
	m.ctx.Mux.Lock()
	defer m.ctx.Mux.Unlock()

	c, cancel := context.WithCancel(context.Background())
	m.consumerPool[r.RuleId] = cancel

	ticker := time.Tick(time.Second)
	go func(ctx context.Context, r models.ProbingRule) {
		for {
			select {
			case <-ticker:
				result, err := m.ctx.Redis.Event().GetPECache(r.GetFiringAlertCacheKey())
				if err == nil {
					m.handleAlert(result)
				}
			case <-ctx.Done():
				return
			}
		}
	}(c, r)
}

func (m *ConsumeProbing) Stop(id string) {
	m.ctx.Mux.Lock()
	defer m.ctx.Mux.Unlock()

	if cancel, exists := m.consumerPool[id]; exists {
		cancel()
	}
}

// 推送告警
func (m *ConsumeProbing) handleAlert(alert models.ProbingEvent) {
	if alert.RuleId == "" {
		return
	}

	if m.filterEvent(alert) {
		r := models.NoticeQuery{
			TenantId: alert.TenantId,
			Uuid:     alert.NoticeId,
		}
		noticeData, _ := ctx.DB.Notice().Get(r)
		alert.DutyUser = process.GetDutyUser(m.ctx, noticeData)

		var content string
		if noticeData.NoticeType == "CustomHook" {
			content = tools.JsonMarshal(alert)
		} else {
			content = templates.NewTemplate(m.ctx, buildEvent(alert), noticeData).CardContentMsg
		}
		err := sender.Sender(m.ctx, sender.SendParams{
			TenantId:    alert.TenantId,
			Severity:    alert.Severity,
			NoticeType:  noticeData.NoticeType,
			NoticeId:    noticeData.Uuid,
			NoticeName:  noticeData.Name,
			IsRecovered: alert.IsRecovered,
			Hook:        noticeData.Hook,
			Email:       noticeData.Email,
			Content:     content,
			Event:       nil,
		})
		if err != nil {
			logc.Errorf(ctx.Ctx, err.Error())
			return
		}
	}
}

func (m *ConsumeProbing) filterEvent(alert models.ProbingEvent) bool {
	var pass bool
	if !alert.IsRecovered {
		if alert.LastSendTime == 0 || alert.LastEvalTime >= alert.LastSendTime+alert.RepeatNoticeInterval*60 {
			alert.LastSendTime = time.Now().Unix()
			m.ctx.Redis.Event().SetPECache(alert, 0)
			return true
		}
	} else {
		removeAlertFromCache(alert)
		return true
	}
	return pass
}

// 删除缓存
func removeAlertFromCache(alert models.ProbingEvent) {
	ctx.DO().Redis.Redis().Del(alert.GetFiringAlertCacheKey())
}

func buildEvent(event models.ProbingEvent) models.AlertCurEvent {
	return models.AlertCurEvent{
		TenantId:               event.TenantId,
		RuleId:                 event.RuleId,
		Fingerprint:            event.Fingerprint,
		Severity:               event.Severity,
		Metric:                 event.Metric,
		NoticeId:               event.NoticeId,
		Annotations:            event.Annotations,
		IsRecovered:            event.IsRecovered,
		FirstTriggerTime:       event.FirstTriggerTime,
		FirstTriggerTimeFormat: event.FirstTriggerTimeFormat,
		RepeatNoticeInterval:   event.RepeatNoticeInterval,
		LastEvalTime:           event.LastEvalTime,
		LastSendTime:           event.LastSendTime,
		RecoverTime:            event.RecoverTime,
		RecoverTimeFormat:      event.RecoverTimeFormat,
		DutyUser:               event.DutyUser,
		RecoverNotify:          event.GetRecoverNotify(),
	}
}
