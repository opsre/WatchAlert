package silence

import (
	"context"
	"sync"
	"time"

	"watchAlert/alert/mute"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"

	"github.com/zeromicro/go-zero/core/logc"
)

const (
	// SilenceRuleProcessInterval controls how often silence rule statuses are refreshed.
	SilenceRuleProcessInterval = time.Second
)

type silenceCache interface {
	GetAlertMutes(tenantId string) ([]string, error)
	WithIdGetMuteFromCache(tenantId, id string) (*models.AlertSilences, error)
	PushAlertMute(mute models.AlertSilences)
}

type silenceRepo interface {
	Update(r models.AlertSilences) error
}

var processOnce sync.Once

// NewSilenceEngine starts the background job that updates silence rule statuses.
// The first scan runs asynchronously, so application initialization is not blocked.
func NewSilenceEngine(appCtx *ctx.Context) {
	if appCtx == nil || appCtx.Ctx == nil {
		return
	}

	processOnce.Do(func() {
		go runSilenceRuleProcessor(appCtx)
	})
}

func runSilenceRuleProcessor(appCtx *ctx.Context) {
	ticker := time.NewTicker(SilenceRuleProcessInterval)
	defer ticker.Stop()

	processAllTenants(appCtx)

	for {
		select {
		case <-ticker.C:
			processAllTenants(appCtx)
		case <-appCtx.Ctx.Done():
			logc.Info(appCtx.Ctx, "静默规则后台任务已停止")
			return
		}
	}
}

func processAllTenants(appCtx *ctx.Context) {
	tenants, err := appCtx.DB.Tenant().List("")
	if err != nil {
		logc.Errorf(appCtx.Ctx, "获取所有租户失败, err: %s", err.Error())
		return
	}

	currentTime := time.Now().Unix()
	for _, tenant := range tenants {
		if tenant.ID == "" {
			logc.Error(appCtx.Ctx, "跳过没有 ID 的租户")
			continue
		}

		process(appCtx.Ctx, appCtx.Redis.Silence(), appCtx.DB.Silence(), tenant.ID, currentTime)
	}
}

func process(logCtx context.Context, silenceCtx silenceCache, silenceRepo silenceRepo, tenantId string, currentTime int64) {
	silenceIds, err := silenceCtx.GetAlertMutes(tenantId)
	if err != nil {
		logc.Errorf(logCtx, "获取租户 %s 的静默规则失败, err: %s", tenantId, err.Error())
		return
	}

	for _, silenceId := range silenceIds {
		muteRule, err := silenceCtx.WithIdGetMuteFromCache(tenantId, silenceId)
		if err != nil {
			logc.Errorf(logCtx, "获取租户 %s 静默规则 %s 失败, err: %s", tenantId, silenceId, err.Error())
			continue
		}
		if muteRule == nil {
			logc.Errorf(logCtx, "租户 %s 静默规则 %s 为空", tenantId, silenceId)
			continue
		}

		newStatus := muteRule.Status
		switch {
		case muteRule.Status != 2 && muteRule.EndsAt <= currentTime:
			newStatus = 2
		case muteRule.Status == 0 && muteRule.StartsAt <= currentTime:
			newStatus = 1
		}

		if newStatus == muteRule.Status {
			continue
		}

		muteRule.Status = newStatus
		if err := silenceRepo.Update(*muteRule); err != nil {
			logc.Errorf(logCtx, "更新租户 %s 静默规则 %s 失败, err: %s", tenantId, silenceId, err.Error())
			continue
		}

		silenceCtx.PushAlertMute(*muteRule)
		mute.InvalidateSilenceCache(tenantId)
	}
}
