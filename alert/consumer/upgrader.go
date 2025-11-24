package consumer

import (
	"fmt"
	"time"
	"watchAlert/alert/mute"
	"watchAlert/alert/process"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"

	"github.com/zeromicro/go-zero/core/logc"
)

// AggregatedAlert 存储聚合后的告警信息
type AggregatedAlert struct {
	Fingerprints []string
	Events       []*models.AlertCurEvent
	Status       int64
	Timeout      int64
}

// alarmUpgrade 处理告警升级主入口
func alarmUpgrade(ctx *ctx.Context, faultCenter models.FaultCenter, alerts map[string]*models.AlertCurEvent) error {
	currentTime := time.Now().Unix()
	if !faultCenter.GetIsUpgradeEnabled() {
		return nil
	}

	// 过滤告警事件
	filterAlerts := filterAlertEvents(faultCenter, alerts)
	if len(filterAlerts) == 0 {
		return nil
	}

	confirmAggregated := createAggregatedAlert(models.ConfirmStatus, faultCenter)
	// 遍历事件并处理升级阶段
	for _, event := range filterAlerts {
		// 确认阶段
		if !event.ConfirmState.IsOk {
			if err := processStage(ctx, faultCenter, event, currentTime, confirmAggregated, models.ConfirmStatus); err != nil {
				logc.Error(ctx.Ctx, fmt.Errorf("process confirm stage failed: %w", err))
			}
		}
	}

	if confirmAggregated != nil {
		sendIfNotEmpty(ctx, faultCenter, confirmAggregated)
	}

	return nil
}

// filterAlertEvents 过滤告警事件
func filterAlertEvents(faultCenter models.FaultCenter, alerts map[string]*models.AlertCurEvent) []*models.AlertCurEvent {
	newEvents := make([]*models.AlertCurEvent, 0, len(alerts))

	for _, event := range alerts {
		switch event.Status {
		case models.StatePreAlert, models.StatePendingRecovery:
			continue
		}

		// 过滤不满足严重级别评估的事件
		if !faultCenter.GetSeverityAssessmentResult(event.Severity) {
			continue
		}

		// 统一过滤恢复状态的事件
		if event.IsRecovered {
			continue
		}

		// 过滤被静默的事件
		if isMutedEvent(event, faultCenter) {
			continue
		}

		newEvents = append(newEvents, event)
	}

	return newEvents
}

// isMutedEvent 检查事件是否被静默
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

// createAggregatedAlert 创建聚合告警对象
func createAggregatedAlert(status int64, faultCenter models.FaultCenter) *AggregatedAlert {
	return &AggregatedAlert{
		Fingerprints: make([]string, 0),
		Events:       make([]*models.AlertCurEvent, 0),
		Status:       status,
		Timeout:      faultCenter.GetTimeout(),
	}
}

// processStage 统一处理确认和处理阶段的逻辑
func processStage(ctx *ctx.Context, faultCenter models.FaultCenter, alert *models.AlertCurEvent, currentTime int64, aggregated *AggregatedAlert, status int64) error {
	// 检查是否超时 (达到升级条件)
	statusStrategy := faultCenter.UpgradeStrategy
	startTime := getStartTime(alert)
	timeoutDuration := statusStrategy.Timeout

	isUpgradeTimeout, err := checkTimeout(startTime, currentTime, timeoutDuration)
	if err != nil {
		return fmt.Errorf("check confirm upgrade timeout failed: %w", err)
	}
	if !isUpgradeTimeout {
		return nil
	}

	// 检查是否需要重新通知 (达到通知间隔)
	lastNoticeTime := getLastNoticeTime(alert)
	if lastNoticeTime != 0 {
		reNotifyDuration := statusStrategy.RepeatInterval

		// 检查通知间隔是否超时
		isReNotifyTimeout, err := checkTimeout(lastNoticeTime, currentTime, reNotifyDuration)
		if err != nil {
			return fmt.Errorf("check confirm re-notify interval failed: %w", err)
		}
		if !isReNotifyTimeout {
			return nil
		}
	}

	aggregated.Fingerprints = append(aggregated.Fingerprints, alert.Fingerprint)
	aggregated.Events = append(aggregated.Events, alert)
	// 更新通知时间，并推送到 Redis
	setLastNoticeTime(ctx, alert, currentTime)

	return nil
}

// getStartTime 根据状态获取开始时间
func getStartTime(alert *models.AlertCurEvent) int64 {
	return alert.FirstTriggerTime
}

// getLastNoticeTime 根据状态获取上次通知时间
func getLastNoticeTime(alert *models.AlertCurEvent) int64 {
	return alert.ConfirmState.ConfirmTimeoutSendTime
}

// setLastNoticeTime 根据状态设置上次通知时间
func setLastNoticeTime(ctx *ctx.Context, alert *models.AlertCurEvent, currentTime int64) {
	alert.ConfirmState.ConfirmTimeoutSendTime = currentTime
	ctx.Redis.Alert().PushAlertEvent(alert)
}

// sendIfNotEmpty 检查聚合告警是否不为空，如果不为空则发送
func sendIfNotEmpty(ctx *ctx.Context, faultCenter models.FaultCenter, aggregated *AggregatedAlert) {
	if len(aggregated.Events) == 0 {
		return
	}

	// 仅保留第一个事件发送
	if len(aggregated.Events) > 1 {
		aggregated.Events[0].Annotations = fmt.Sprintf("%s%s", aggregated.Events[0].Annotations, getContent(len(aggregated.Events)))
		aggregated.Events = aggregated.Events[:1]
	}

	if err := sendAggregatedAlert(ctx, faultCenter, aggregated); err != nil {
		logc.Error(ctx.Ctx, fmt.Errorf("send aggregated confirm alert failed: %w", err))
	}
}

// sendAggregatedAlert 发送聚合后的告警函数
func sendAggregatedAlert(ctx *ctx.Context, faultCenter models.FaultCenter, aggregated *AggregatedAlert) error {
	noticeId := faultCenter.GetUpgradeNoticeId()
	if noticeId == "" {
		return nil
	}

	logc.Alert(ctx.Ctx, fmt.Sprintf("Aggregated alarm confirm timeout fingerprints: %v, exceeded %d min",
		aggregated.Fingerprints,
		aggregated.Timeout))

	return process.HandleAlert(ctx, "upgrade", faultCenter, noticeId, aggregated.Events)
}

// getContent 生成聚合通知内容
func getContent(number int) string {
	return fmt.Sprintf("\n聚合 %d 条升级通知, 详情请前往 WatchAlert 查看", number)
}

// checkTimeout 检查是否超时 duration 单位为分钟
func checkTimeout(startTime, currentTime int64, duration int64) (bool, error) {
	timeoutSeconds := duration * 60

	// 如果 currentTime 超过 startTime + timeoutSeconds，则超时
	return currentTime > startTime+timeoutSeconds, nil
}
