package mute

import (
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"
)

type MuteParams struct {
	RecoverNotify *bool
	IsRecovered   bool
	TenantId      string
	Labels        map[string]interface{}
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
	if mp.RecoverNotify == nil {
		return mp.IsRecovered
	}
	return mp.IsRecovered && !*mp.RecoverNotify
}

// IsSilence 判断是否静默
// 静默规则列表通过进程内缓存读取（getTenantSilences），
// 避免高事件量场景下对每个事件重复发起 Redis 请求。
// 规则变更后通过 InvalidateSilenceCache 主动失效。
func IsSilence(mute MuteParams) bool {
	for _, muteRule := range getTenantSilences(mute.TenantId) {
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
			matched = (val == muteLabel.Value)
		case "!=":
			matched = (val != muteLabel.Value)
		case "=~":
			matched, _ = tools.MatchRegexp(muteLabel.Value, val)
		case "!~":
			if regexpMatched, err := tools.MatchRegexp(muteLabel.Value, val); err == nil {
				matched = !regexpMatched
			}
		default:
			matched = false
		}

		if !matched {
			return false // 只要有一个不匹配，就不静默
		}
	}

	return true
}
