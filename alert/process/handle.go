package process

import (
	"fmt"
	"strings"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/pkg/sender"
	"watchAlert/pkg/templates"
	"watchAlert/pkg/tools"

	"github.com/zeromicro/go-zero/core/logc"
	"golang.org/x/sync/errgroup"
)

// HandleAlert 处理告警逻辑
func HandleAlert(ctx *ctx.Context, faultCenter models.FaultCenter, noticeId string, alerts []*models.AlertCurEvent) error {
	curTime := time.Now().Unix()
	g := new(errgroup.Group)

	// 获取通知对象详细信息
	noticeData, err := getNoticeData(ctx, faultCenter.TenantId, noticeId)
	if err != nil {
		logc.Error(ctx.Ctx, fmt.Sprintf("Failed to get notice data: %v", err))
		return err
	}

	// 按告警等级分组
	severityGroups := make(map[string][]*models.AlertCurEvent)
	for _, alert := range alerts {
		severityGroups[alert.Severity] = append(severityGroups[alert.Severity], alert)
	}

	// 告警聚合
	aggregationEvents := alarmAggregation(ctx, faultCenter, severityGroups)
	for severity, events := range aggregationEvents {
		g.Go(func() error {
			if events == nil {
				return nil
			}

			// 获取当前事件等级对应的 Hook 和 Sign
			Hook, Sign := getNoticeHookUrlAndSign(noticeData, severity)

			for _, event := range events {
				if !event.IsRecovered {
					event.LastSendTime = curTime
					ctx.Redis.Alert().PushAlertEvent(event)
				}

				phoneNumber := func() []string {
					if len(event.DutyUserPhoneNumber) > 0 {
						return event.DutyUserPhoneNumber
					}
					if len(noticeData.PhoneNumber) > 0 {
						return noticeData.PhoneNumber
					}
					return []string{}
				}()

				event.DutyUser = strings.Join(GetDutyUsers(ctx, noticeData), " ")
				event.DutyUserPhoneNumber = GetDutyUserPhoneNumber(ctx, noticeData)
				content := generateAlertContent(ctx, event, noticeData)
				err := sender.Sender(ctx, sender.SendParams{
					TenantId:    event.TenantId,
					EventId:     event.EventId,
					RuleName:    event.RuleName,
					Severity:    event.Severity,
					NoticeType:  noticeData.NoticeType,
					NoticeId:    noticeId,
					NoticeName:  noticeData.Name,
					IsRecovered: event.IsRecovered,
					Hook:        Hook,
					Email:       getNoticeEmail(noticeData, severity),
					Content:     content,
					PhoneNumber: phoneNumber,
					Sign:        Sign,
				})
				if err != nil {
					logc.Error(ctx.Ctx, fmt.Sprintf("Failed to send alert: %v", err))
				}
			}
			return nil
		})
	}

	return g.Wait()
}

// alarmAggregation 告警聚合
func alarmAggregation(ctx *ctx.Context, faultCenter models.FaultCenter, alertGroups map[string][]*models.AlertCurEvent) map[string][]*models.AlertCurEvent {
	curTime := time.Now().Unix()
	newAlertGroups := alertGroups
	switch faultCenter.GetAlarmAggregationType() {
	case "Rule":
		for severity, events := range alertGroups {
			newAlertGroups[severity] = withRuleGroupByAlerts(ctx, curTime, events)
		}
	default:
		return alertGroups
	}

	return newAlertGroups
}

// withRuleGroupByAlerts 聚合告警
func withRuleGroupByAlerts(ctx *ctx.Context, timeInt int64, alerts []*models.AlertCurEvent) []*models.AlertCurEvent {
	if len(alerts) <= 1 {
		return alerts
	}

	var aggregatedAlert *models.AlertCurEvent
	for i := range alerts {
		alert := alerts[i]
		if !strings.Contains(alert.Annotations, "聚合") {
			alert.Annotations += fmt.Sprintf("\n聚合 %d 条告警\n", len(alerts))
		}
		aggregatedAlert = alert

		if !alert.IsRecovered {
			alert.LastSendTime = timeInt
			ctx.Redis.Alert().PushAlertEvent(alert)
		}
	}

	return []*models.AlertCurEvent{aggregatedAlert}
}

// getNoticeData 获取 Notice 数据
func getNoticeData(ctx *ctx.Context, tenantId, noticeId string) (models.AlertNotice, error) {
	return ctx.DB.Notice().Get(tenantId, noticeId)
}

// getNoticeHookUrlAndSign 获取事件等级对应的 Hook 和 Sign
func getNoticeHookUrlAndSign(notice models.AlertNotice, severity string) (string, string) {
	if notice.Routes != nil {
		for _, hook := range notice.Routes {
			if hook.Severity == severity {
				return hook.Hook, hook.Sign
			}
		}
	}
	return notice.DefaultHook, notice.DefaultSign
}

// getNoticeEmail 获取事件等级对应的 Email
func getNoticeEmail(notice models.AlertNotice, severity string) models.Email {
	if notice.Routes != nil {
		for _, route := range notice.Routes {
			if route.Severity == severity {
				return models.Email{
					Subject: notice.Email.Subject,
					To:      route.To,
					CC:      route.CC,
				}
			}
		}
	}
	return notice.Email
}

// generateAlertContent 生成告警内容
func generateAlertContent(ctx *ctx.Context, alert *models.AlertCurEvent, noticeData models.AlertNotice) string {
	if noticeData.NoticeType == "CustomHook" {
		return tools.JsonMarshalToString(alert)
	}
	return templates.NewTemplate(ctx, *alert, noticeData).CardContentMsg
}
