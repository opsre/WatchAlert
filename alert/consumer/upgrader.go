package consumer

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"time"
	"watchAlert/alert/mute"
	"watchAlert/alert/process"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
)

// alarmUpgrade 处理告警升级主入口
func alarmUpgrade(ctx *ctx.Context, faultCenter models.FaultCenter, alerts map[string]*models.AlertCurEvent) error {
	if !faultCenter.GetIsUpgradeEnabled() {
		return nil
	}

	filterAlerts := filterAlertEvents(ctx.Ctx, faultCenter, alerts)
	currentTime := time.Now().Unix()

	confirmAggregated := createAggregatedAlert(models.ConfirmStatus, faultCenter)
	handleAggregated := createAggregatedAlert(models.HandleStatus, faultCenter)

	for _, event := range filterAlerts {
		if !event.UpgradeState.IsConfirm && faultCenter.GetStrategy(models.ConfirmStatus).GetEnabled() {
			if err := processStage(ctx, faultCenter, event, currentTime, confirmAggregated, models.ConfirmStatus); err != nil {
				return fmt.Errorf("process confirm stage failed: %w", err)
			}
		}

		if event.UpgradeState.IsConfirm && !event.UpgradeState.IsHandle && faultCenter.GetStrategy(models.HandleStatus).GetEnabled() {
			if err := processStage(ctx, faultCenter, event, currentTime, handleAggregated, models.HandleStatus); err != nil {
				return fmt.Errorf("process handle stage failed: %w", err)
			}
		}
	}

	sendIfNotEmpty(ctx, faultCenter, confirmAggregated)
	sendIfNotEmpty(ctx, faultCenter, handleAggregated)

	return nil
}

// createAggregatedAlert 创建聚合告警对象
func createAggregatedAlert(status int64, faultCenter models.FaultCenter) *AggregatedAlert {
	return &AggregatedAlert{
		Status:  status,
		Timeout: faultCenter.GetTimeout(status),
	}
}

// processStage 统一处理确认和处理阶段的逻辑
func processStage(ctx *ctx.Context, faultCenter models.FaultCenter, alert *models.AlertCurEvent, currentTime int64, aggregated *AggregatedAlert, status int64) error {
	startTime := getStartTime(alert, status)
	timeout, err := checkTimeout(startTime, currentTime, faultCenter.GetTimeout(status))
	if err != nil {
		return fmt.Errorf("check %s timeout failed: %w", getStatusString(status), err)
	}

	if !timeout {
		return nil
	}

	// 检查是否需要重新通知
	if lastNoticeTime := getLastNoticeTime(alert, status); lastNoticeTime != 0 {
		reNotify, err := checkTimeout(lastNoticeTime, currentTime, faultCenter.GetNoticeInterval(status))
		if err != nil {
			return fmt.Errorf("check %s re-notify timeout failed: %w", getStatusString(status), err)
		}
		if !reNotify {
			return nil
		}
	}

	// TODO 区分告警 / 恢复事件
	aggregated.Fingerprints = append(aggregated.Fingerprints, alert.Fingerprint)
	aggregated.Events = append(aggregated.Events, alert)
	setLastNoticeTime(alert, status, currentTime)
	processAlarmEvent(ctx, status, *alert, currentTime)

	return nil
}

// getStartTime 根据状态获取开始时间
func getStartTime(alert *models.AlertCurEvent, status int64) int64 {
	switch status {
	case models.ConfirmStatus:
		return alert.FirstTriggerTime
	case models.HandleStatus:
		return alert.UpgradeState.ConfirmOkTime
	default:
		return 0
	}
}

// getLastNoticeTime 根据状态获取上次通知时间
func getLastNoticeTime(alert *models.AlertCurEvent, status int64) int64 {
	switch status {
	case models.ConfirmStatus:
		return alert.UpgradeState.ConfirmSendTime
	case models.HandleStatus:
		return alert.UpgradeState.HandleSendTime
	default:
		return 0
	}
}

// setLastNoticeTime 根据状态设置上次通知时间
func setLastNoticeTime(alert *models.AlertCurEvent, status int64, currentTime int64) {
	switch status {
	case models.ConfirmStatus:
		alert.UpgradeState.ConfirmSendTime = currentTime
	case models.HandleStatus:
		alert.UpgradeState.HandleSendTime = currentTime
	}
}

// sendIfNotEmpty 检查聚合告警是否不为空，如果不为空则发送
func sendIfNotEmpty(ctx *ctx.Context, faultCenter models.FaultCenter, aggregated *AggregatedAlert) {
	if len(aggregated.Events) == 0 {
		return
	}

	if len(aggregated.Events) > 1 {
		aggregated.Events[0].Annotations = aggregated.Events[0].Annotations + getContent(len(aggregated.Events))
		aggregated.Events = []*models.AlertCurEvent{aggregated.Events[0]}
	}

	if err := sendAggregatedAlert(ctx, faultCenter, aggregated); err != nil {
		logc.Error(ctx.Ctx, fmt.Errorf("send aggregated %s alert failed: %w", getStatusString(aggregated.Status), err))
	}
}

// getContent 生成聚合通知内容
func getContent(number int) string {
	return fmt.Sprintf("\n聚合 %d 条升级通知, 详情请前往 WatchAlert 查看", number)
}

// filterAlertEvents 过滤告警事件
func filterAlertEvents(ctx context.Context, faultCenter models.FaultCenter, alerts map[string]*models.AlertCurEvent) []*models.AlertCurEvent {
	var newEvents []*models.AlertCurEvent
	for _, event := range alerts {
		// 过滤掉 预告警, 待恢复 状态的事件
		if event.Status == models.StatePreAlert || event.Status == models.StatePendingRecovery || event.IsRecovered {
			continue
		}

		if !faultCenter.GetSeverityAssessmentResult(event.Severity) {
			continue
		}

		if isMutedEvent(event, faultCenter) {
			continue
		}

		newEvents = append(newEvents, event)
	}

	return newEvents
}

func isMutedEvent(event *models.AlertCurEvent, faultCenter models.FaultCenter) bool {
	return mute.IsMuted(mute.MuteParams{
		EffectiveTime: event.EffectiveTime,
		IsRecovered:   event.IsRecovered,
		TenantId:      event.TenantId,
		Labels:        event.Labels,
		FaultCenterId: event.FaultCenterId,
		RecoverNotify: faultCenter.RecoverNotify,
	})
}

// sendAggregatedAlert 发送聚合后的告警函数
func sendAggregatedAlert(ctx *ctx.Context, faultCenter models.FaultCenter, aggregated *AggregatedAlert) error {
	var noticeId string
	if aggregated.Status == models.ConfirmStatus {
		noticeId = faultCenter.GetUpgradeNoticeId(models.ConfirmStatus)
	} else if aggregated.Status == models.HandleStatus {
		noticeId = faultCenter.GetUpgradeNoticeId(models.HandleStatus)
	}

	logc.Alert(ctx.Ctx, fmt.Sprintf("Aggregated alarm %s timeout fingerprints: %v, exceeded %d min",
		getStatusString(aggregated.Status),
		aggregated.Fingerprints,
		aggregated.Timeout))

	err := process.HandleAlert(ctx, faultCenter, noticeId, aggregated.Events)
	if err != nil {
		return err
	}

	return nil
}

// getStatusString 辅助函数，用于获取状态的字符串表示
func getStatusString(status int64) string {
	switch status {
	case models.ConfirmStatus:
		return "confirm"
	case models.HandleStatus:
		return "handle"
	default:
		return "unknown"
	}
}

// checkTimeout 检查是否超时
func checkTimeout(startTime, currentTime int64, duration int64) (bool, error) {
	timeoutDuration := time.Duration(duration) * time.Minute
	return currentTime > time.Unix(startTime, 0).Add(timeoutDuration).Unix(), nil
}

func processAlarmEvent(ctx *ctx.Context, status int64, alertEvent models.AlertCurEvent, curT int64) {
	cache := ctx.Redis.Alert()
	event, err := cache.GetEventFromCache(alertEvent.TenantId, alertEvent.FaultCenterId, alertEvent.Fingerprint)
	if err != nil {
		logc.Error(ctx.Ctx, fmt.Sprintf("get event info fail, err: %s", err.Error()))
		return
	}

	switch status {
	case models.ConfirmStatus:
		event.UpgradeState.ConfirmSendTime = curT
	case models.HandleStatus:
		event.UpgradeState.HandleSendTime = curT
	}

	cache.PushAlertEvent(&event)
}

// AggregatedAlert 存储聚合后的告警信息
type AggregatedAlert struct {
	Fingerprints []string
	Events       []*models.AlertCurEvent
	Status       int64
	Timeout      int64
}
