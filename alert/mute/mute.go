package mute

import (
	"regexp"
	"watchAlert/internal/ctx"
	models "watchAlert/internal/models"

	"github.com/zeromicro/go-zero/core/logc"
)

type MuteParams struct {
	RecoverNotify *bool
	IsRecovered   bool
	TenantId      string
	Labels        map[string]interface{}
	FaultCenterId string
}

func IsMuted(mute MuteParams) bool {
	if IsSilence(mute) {
		return true
	}

	if RecoverNotify(mute) {
		return true
	}

	return false
}

// RecoverNotify 判断是否推送恢复通知
func RecoverNotify(mp MuteParams) bool {
	return mp.IsRecovered && !*mp.RecoverNotify
}

// IsSilence 判断是否静默
func IsSilence(mute MuteParams) bool {
	silenceCtx := ctx.Redis.Silence()
	// 获取静默列表中所有的id
	ids, err := silenceCtx.GetAlertMutes(mute.TenantId, mute.FaultCenterId)
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

		if evalCondition(mute.Labels, muteRule.Labels) {
			return true
		}
	}

	return false
}

func evalCondition(metrics map[string]interface{}, muteLabels []models.SilenceLabel) bool {
	for _, muteLabel := range muteLabels {
		value, exists := metrics[muteLabel.Key]
		if !exists {
			return false
		}

		val, ok := value.(string)
		if !ok {
			continue
		}

		var matched bool
		switch muteLabel.Operator {
		case "==", "=":
			matched = regexp.MustCompile(muteLabel.Value).MatchString(val)
		case "!=":
			matched = !regexp.MustCompile(muteLabel.Value).MatchString(val)
		default:
			matched = false
		}

		if !matched {
			return false // 只要有一个不匹配，就不静默
		}
	}

	return true
}
