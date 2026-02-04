package consumer

import (
	"fmt"
	"slices"
	"strings"
	"time"
	"watchAlert/alert/mute"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/pkg/sender"
	"watchAlert/pkg/templates"
	"watchAlert/pkg/tools"

	"github.com/zeromicro/go-zero/core/logc"
	"golang.org/x/sync/errgroup"
)

// handleAlert 处理告警逻辑
func handleAlert(ctx *ctx.Context, processType string, faultCenter models.FaultCenter, noticeId string, alerts []*models.AlertCurEvent) error {
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
	var aggregationEvents map[string][]*models.AlertCurEvent
	if processType == "alarm" {
		aggregationEvents = alarmAggregation(ctx, processType, faultCenter, severityGroups)
	} else {
		aggregationEvents = severityGroups
	}

	for severity, events := range aggregationEvents {
		g.Go(func() error {
			if events == nil {
				return nil
			}

			// 获取当前事件等级对应的路由配置
			routes := getNoticeRoutes(noticeData, severity)
			for _, event := range events {
				if event.Fingerprint == "" {
					continue
				}

				if processType == "alarm" && !event.IsRecovered {
					event.LastSendTime = curTime
					ctx.Redis.Alert().PushAlertEvent(event)
				}

				if len(routes) == 0 {
					logc.Infof(ctx.Ctx, "没用匹配的通知策略, 告警事件名称: %s, 通知对象名称: %s", event.RuleName, noticeData.Name)
				}

				if mute.IsMuted(mute.MuteParams{
					IsRecovered:   event.IsRecovered,
					TenantId:      event.TenantId,
					Labels:        event.Labels,
					FaultCenterId: event.FaultCenterId,
					RecoverNotify: faultCenter.RecoverNotify,
				}) {
					continue
				}

				for _, route := range routes {
					// 设置值班用户信息
					event.DutyUser = strings.Join(getDutyUsers(ctx, noticeData, route.NoticeType), " ")

					// 生成告警内容
					content := generateAlertContent(ctx, event, noticeData, route)

					// 构建邮件信息
					email := models.Email{
						Subject: route.Subject,
						To:      route.To,
						CC:      route.CC,
					}

					// 发送告警
					err := sender.Sender(ctx, sender.SendParams{
						TenantId:    event.TenantId,
						EventId:     event.EventId,
						RuleName:    event.RuleName,
						Severity:    event.Severity,
						NoticeType:  route.NoticeType,
						NoticeId:    noticeId,
						NoticeName:  noticeData.Name,
						IsRecovered: event.IsRecovered,
						Hook:        route.Hook,
						Email:       email,
						Content:     content,
						Sign:        route.Sign,
					})
					if err != nil {
						logc.Error(ctx.Ctx, fmt.Sprintf("Failed to send alert: %v", err))
					}
				}
			}

			return nil
		})
	}

	return g.Wait()
}

// alarmAggregation 告警聚合
func alarmAggregation(ctx *ctx.Context, processType string, faultCenter models.FaultCenter, alertGroups map[string][]*models.AlertCurEvent) map[string][]*models.AlertCurEvent {
	// 仅当 processType 为 "alarm" 时执行聚合
	if processType != "alarm" {
		return alertGroups
	}

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

	for i := range alerts {
		alert := alerts[i]
		if !alert.IsRecovered {
			alert.LastSendTime = timeInt
			ctx.Redis.Alert().PushAlertEvent(alert)
		}
	}

	event := *alerts[0]
	return []*models.AlertCurEvent{&event}
}

// getNoticeData 获取 Notice 数据
func getNoticeData(ctx *ctx.Context, tenantId, noticeId string) (models.AlertNotice, error) {
	return ctx.DB.Notice().Get(tenantId, noticeId)
}

// getNoticeRoutes 获取事件等级对应的路由配置
func getNoticeRoutes(notice models.AlertNotice, severity string) []models.Route {
	var routes []models.Route
	if notice.Routes != nil {
		for _, route := range notice.Routes {
			if slices.Contains(route.Severitys, severity) {
				routes = append(routes, route)
			}
		}
	}

	return routes
}

type WebhookContent struct {
	Alarm     *models.AlertCurEvent `json:"alarm"`
	DutyUsers []models.DutyUser     `json:"dutyUsers"`
}

// generateAlertContent 生成告警内容
func generateAlertContent(ctx *ctx.Context, alert *models.AlertCurEvent, noticeData models.AlertNotice, route models.Route) string {
	if route.NoticeType == "WebHook" {
		users, ok := ctx.DB.DutyCalendar().GetDutyUserInfo(*noticeData.GetDutyId(), time.Now().Format("2006-1-2"))
		if !ok || len(users) == 0 {
			logc.Error(ctx.Ctx, "Failed to get duty users, noticeName: ", noticeData.Name)
		}

		var dutyUsers = []models.DutyUser{}
		for _, user := range users {
			dutyUsers = append(dutyUsers, models.DutyUser{
				Email:    user.Email,
				Mobile:   user.Phone,
				UserId:   user.UserId,
				Username: user.UserName,
			})
		}
		content := WebhookContent{
			Alarm:     alert,
			DutyUsers: dutyUsers,
		}

		return tools.JsonMarshalToString(content)
	}

	template, err := templates.NewTemplate(ctx, *alert, route)
	if err != nil {
		logc.Error(ctx.Ctx, fmt.Sprintf("Failed to create template: %v", err))
		return ""
	}
	return template.CardContentMsg
}

func getDutyUsers(ctx *ctx.Context, noticeData models.AlertNotice, noticeType string) []string {
	var us []string
	users, ok := ctx.DB.DutyCalendar().GetDutyUserInfo(*noticeData.GetDutyId(), time.Now().Format("2006-1-2"))
	if ok {
		switch noticeType {
		case "FeiShu":
			for _, user := range users {
				us = append(us, fmt.Sprintf("<at id=%s></at>", user.DutyUserId))
			}
			return us
		case "DingDing":
			for _, user := range users {
				us = append(us, fmt.Sprintf("@%s", user.DutyUserId))
			}
			return us
		case "Email", "WeChat", "WebHook":
			for _, user := range users {
				us = append(us, fmt.Sprintf("@%s", user.UserName))
			}
			return us
		case "Slack":
			for _, user := range users {
				us = append(us, fmt.Sprintf("<@%s>", user.DutyUserId))
			}
			return us
		}
	}

	return []string{"暂无"}
}
