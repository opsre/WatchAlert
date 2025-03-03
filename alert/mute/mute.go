package mute

import (
	"github.com/zeromicro/go-zero/core/logc"
	"time"
	models "watchAlert/internal/models"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/tools"
)

type MuteParams struct {
	EffectiveTime models.EffectiveTime
	RecoverNotify *bool
	IsRecovered   bool
	TenantId      string
	Metrics       map[string]interface{}
	FaultCenterId string
}

func IsMuted(mute MuteParams) bool {
	if IsSilence(mute) {
		return true
	}

	if NotInTheEffectiveTime(mute) {
		return true
	}

	if RecoverNotify(mute) {
		return true
	}

	return false
}

// NotInTheEffectiveTime 判断生效时间
func NotInTheEffectiveTime(mp MuteParams) bool {
	if len(mp.EffectiveTime.Week) <= 0 {
		return false
	}

	// 获取当前日期
	currentTime := time.Now()
	currentWeekday := tools.TimeTransformToWeek(currentTime)
	for _, weekday := range mp.EffectiveTime.Week {
		if currentWeekday == weekday {
			currentTimeSeconds := tools.TimeTransformToSeconds(currentTime)
			return currentTimeSeconds < mp.EffectiveTime.StartTime || currentTimeSeconds > mp.EffectiveTime.EndTime
		}
	}

	return true
}

// RecoverNotify 判断是否推送恢复通知
func RecoverNotify(mp MuteParams) bool {
	return mp.IsRecovered && !*mp.RecoverNotify
}

// IsSilence 判断是否静默
func IsSilence(mute MuteParams) bool {
	silenceCtx := ctx.Redis.Silence()
	// 获取静默列表中所有的id
	ids, err := silenceCtx.GetMutesForFaultCenter(mute.TenantId, mute.FaultCenterId)
	if err != nil {
		logc.Errorf(ctx.Ctx, err.Error())
		return false
	}

	// 根据ID获取到详细的静默规则
	for _, id := range ids {
		muteRule, err := silenceCtx.WithIdGetMuteFromCache(mute.TenantId, mute.FaultCenterId, id)
		if err != nil {
			logc.Errorf(ctx.Ctx, err.Error())
			return false
		}

		if muteRule.Status != 1 {
			continue
		}

		for _, label := range muteRule.Labels {
			if evalCondition(mute.Metrics, label) {
				return true
			}
		}
	}

	return false
}

func evalCondition(metrics map[string]interface{}, c models.SilenceLabel) bool {
	switch c.Operator {
	case "==", "=":
		return metrics[c.Key] == c.Value
	case "!=":
		return metrics[c.Key] != c.Value
	default:
		return false
	}
}
