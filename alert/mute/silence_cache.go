package mute

import (
	"sync"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"

	"github.com/zeromicro/go-zero/core/logc"
)

// 静默规则进程内缓存的过期时间。
// 在 TTL 内对静默规则的判断走本地缓存，避免高事件量场景下
// 每秒对每个事件发起 1+N 次 Redis 往返。
const silenceCacheTTL = 60 * time.Second

type silenceCacheEntry struct {
	rules    []*models.AlertSilences
	expireAt time.Time
}

var (
	silenceCacheMu     sync.RWMutex
	silenceCacheTables = make(map[string]silenceCacheEntry)
)

// getTenantSilences 获取租户的静默规则列表，优先读进程内缓存。
// 缓存不存在或已过期时从 Redis 全量加载一次。
func getTenantSilences(tenantId string) []*models.AlertSilences {
	silenceCacheMu.RLock()
	entry, ok := silenceCacheTables[tenantId]
	silenceCacheMu.RUnlock()

	if ok && time.Now().Before(entry.expireAt) {
		return entry.rules
	}

	rules := loadTenantSilences(tenantId)
	if rules == nil {
		// 加载失败时保留旧缓存（若存在），避免 Redis 抖动导致静默判断全部失效
		if ok {
			return entry.rules
		}
		return nil
	}

	silenceCacheMu.Lock()
	silenceCacheTables[tenantId] = silenceCacheEntry{
		rules:    rules,
		expireAt: time.Now().Add(silenceCacheTTL),
	}
	silenceCacheMu.Unlock()

	return rules
}

// loadTenantSilences 从 Redis 加载租户的全部静默规则。
// 单条规则加载失败时跳过该条并记录日志，不影响其余规则。
func loadTenantSilences(tenantId string) []*models.AlertSilences {
	silenceCtx := ctx.Redis.Silence()

	ids, err := silenceCtx.GetAlertMutes(tenantId)
	if err != nil {
		logc.Errorf(ctx.Ctx, "加载租户静默规则列表失败, tenantId: %s, err: %s", tenantId, err.Error())
		return nil
	}

	rules := make([]*models.AlertSilences, 0, len(ids))
	for _, id := range ids {
		muteRule, err := silenceCtx.WithIdGetMuteFromCache(tenantId, id)
		if err != nil {
			// 静默规则可能已被删除，跳过而非中断整个租户的加载
			logc.Errorf(ctx.Ctx, "获取租户静默规则失败, tenantId: %s, id: %s, err: %s", tenantId, id, err.Error())
			continue
		}
		rules = append(rules, muteRule)
	}

	return rules
}

// InvalidateSilenceCache 使指定租户的静默规则进程内缓存失效。
// 不传参数时清空全部租户的缓存。
// 应在静默规则发生增删改或状态变更后调用，保证判断实时性。
func InvalidateSilenceCache(tenantIds ...string) {
	silenceCacheMu.Lock()
	defer silenceCacheMu.Unlock()

	if len(tenantIds) == 0 {
		silenceCacheTables = make(map[string]silenceCacheEntry)
		return
	}

	for _, tenantId := range tenantIds {
		delete(silenceCacheTables, tenantId)
	}
}
