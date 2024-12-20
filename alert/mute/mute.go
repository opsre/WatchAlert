package mute

import (
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
	Fingerprint   string
}

func IsMuted(mute MuteParams) bool {
	if IsSilence(mute) {
		return true
	}

	if InTheEffectiveTime(mute) {
		return true
	}

	if RecoverNotify(mute) {
		return true
	}

	return false
}

// InTheEffectiveTime 判断生效时间
func InTheEffectiveTime(mp MuteParams) bool {
	if len(mp.EffectiveTime.Week) <= 0 {
		return false
	}

	var (
		p           bool
		currentTime = time.Now()
	)

	cwd := tools.TimeTransformToWeek(currentTime)
	for _, wd := range mp.EffectiveTime.Week {
		if cwd != wd {
			continue
		}
		p = true
	}

	if !p {
		return true
	}

	cts := tools.TimeTransformToSeconds(currentTime)
	if cts < mp.EffectiveTime.StartTime || cts > mp.EffectiveTime.EndTime {
		return true
	}

	return false
}

// RecoverNotify 判断是否推送恢复通知
func RecoverNotify(mp MuteParams) bool {
	// 如果是恢复告警，并且 恢复通知 == 1，即关闭恢复通知
	if mp.IsRecovered && !*mp.RecoverNotify {
		return true
	}

	return false
}

// IsSilence 判断是否静默
func IsSilence(mute MuteParams) bool {
	_, ok := ctx.Redis.Silence().GetCache(models.AlertSilenceQuery{
		TenantId:    mute.TenantId,
		Fingerprint: mute.Fingerprint,
	})

	if ok {
		return true
	} else {
		ttl, _ := ctx.Redis.Redis().TTL(mute.TenantId + ":" + models.SilenceCachePrefix + mute.Fingerprint).Result()
		// 如果剩余生存时间小于0，表示键已过期
		if ttl < 0 {
			// 过期后标记为1
			ctx.DB.DB().Model(models.AlertSilences{}).
				Where("fingerprint = ? and status = ?", mute.Fingerprint, 0).
				Update("status", 1)
		}
	}

	return false
}
