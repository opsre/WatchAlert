package probing

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
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
	go m.Watch(c, r)
}

func (m *ConsumeProbing) Stop(id string) {
	m.ctx.Mux.Lock()
	defer m.ctx.Mux.Unlock()

	if cancel, exists := m.consumerPool[id]; exists {
		cancel()
	}
}

func (m *ConsumeProbing) Watch(ctx context.Context, r models.ProbingRule) {
	taskChan := make(chan struct{}, 1)
	timer := time.NewTicker(time.Second * time.Duration(1))
	defer func() {
		timer.Stop()
	}()

	for {
		select {
		case <-timer.C:
			taskChan <- struct{}{}
			m.executeTask(taskChan, r)
		case <-ctx.Done():
			return
		}
	}
}

func (m *ConsumeProbing) executeTask(taskChan chan struct{}, r models.ProbingRule) {
	defer func() {
		<-taskChan
	}()

	var now = time.Now().Unix()
	event, err := m.ctx.Redis.Probing().GetProbingEventCache(models.BuildProbingEventCacheKey(r.TenantId, r.RuleId))
	if err != nil {
		if err == redis.Nil {
			return
		}
		logc.Error(context.Background(), fmt.Sprintf("获取拨测事件失败, %s", err.Error()))
		return
	}

	if !m.filterEvent(event) {
		return
	}

	newEvent := event
	newEvent.LastSendTime = now
	m.ctx.Redis.Probing().SetProbingEventCache(newEvent, 0)
	m.sendAlert(newEvent)
}

func (m *ConsumeProbing) filterEvent(alert models.ProbingEvent) bool {
	if !alert.IsRecovered {
		if alert.LastSendTime == 0 || alert.LastEvalTime >= alert.LastSendTime+alert.RepeatNoticeInterval*60 {
			return true
		}
	} else {
		m.removeAlertFromCache(alert)
		return true
	}

	return false
}

// 推送告警
func (m *ConsumeProbing) sendAlert(alert models.ProbingEvent) {
	r := models.NoticeQuery{
		TenantId: alert.TenantId,
		Uuid:     alert.NoticeId,
	}
	noticeData, err := m.ctx.DB.Notice().Get(r)
	if err != nil {
		logc.Error(m.ctx.Ctx, "获取通知对象失败, ", err.Error())
		return
	}

	alert.DutyUser = process.GetDutyUser(m.ctx, noticeData)
	err = sender.Sender(m.ctx, sender.SendParams{
		RuleName:    alert.RuleName,
		TenantId:    alert.TenantId,
		NoticeType:  noticeData.NoticeType,
		NoticeId:    noticeData.Uuid,
		NoticeName:  noticeData.Name,
		IsRecovered: alert.IsRecovered,
		Hook:        noticeData.DefaultHook,
		Email:       noticeData.Email,
		Content:     m.getContent(alert, noticeData),
		Sign:        noticeData.DefaultSign,
	})
	if err != nil {
		logc.Errorf(ctx.Ctx, err.Error())
		return
	}
}

func (m *ConsumeProbing) getContent(alert models.ProbingEvent, noticeData models.AlertNotice) string {
	if noticeData.NoticeType == "CustomHook" {
		return tools.JsonMarshal(alert)
	} else {
		return templates.NewTemplate(m.ctx, buildEvent(alert), noticeData).CardContentMsg
	}
}

// 删除缓存
func (m *ConsumeProbing) removeAlertFromCache(alert models.ProbingEvent) {
	m.ctx.Redis.Redis().Del(string(models.BuildProbingEventCacheKey(alert.TenantId, alert.RuleId)))
}

func buildEvent(event models.ProbingEvent) models.AlertCurEvent {
	return models.AlertCurEvent{
		TenantId:               event.TenantId,
		RuleId:                 event.RuleId,
		Fingerprint:            event.Fingerprint,
		Metric:                 event.Metric,
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
	}
}
