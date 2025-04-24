package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	"sync"
	"watchAlert/pkg/tools"
)

type (
	// PendingRecoverCache 用于管理待恢复的告警事件
	PendingRecoverCache struct {
		rc    *redis.Client
		mutex sync.RWMutex
	}

	// PendingRecoverCacheInterface 定义了待恢复的告警事件缓存的操作接口
	PendingRecoverCacheInterface interface {
		Set(tenantId, ruleId, fingerprint string, time int64)
		Get(tenantId, ruleId, fingerprint string) (int64, error)
		Delete(tenantId, ruleId, fingerprint string)
		List(tenantId, ruleId string) map[string]int64
	}

	PendingRecoverCacheKey string

	PendingRecoverCacheData struct {
		Fingerprint string
		Time        int64
	}
)

// newPendingRecoverCacheInterface 创建一个新的 PendingRecoverCache 实例
func newPendingRecoverCacheInterface(r *redis.Client) PendingRecoverCacheInterface {
	return &PendingRecoverCache{
		rc: r,
	}
}

func (p *PendingRecoverCache) Set(tenantId, ruleId, fingerprint string, time int64) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.rc.HSet(string(BuildPendingRecoverCacheKey(tenantId, ruleId)), fingerprint, time)
}

func (p *PendingRecoverCache) Get(tenantId, ruleId, fingerprint string) (int64, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.rc.HGet(string(BuildPendingRecoverCacheKey(tenantId, ruleId)), fingerprint).Int64()
}

func (p *PendingRecoverCache) Delete(tenantId, ruleId, fingerprint string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.rc.HDel(string(BuildPendingRecoverCacheKey(tenantId, ruleId)), fingerprint)
}

func (p *PendingRecoverCache) List(tenantId, ruleId string) map[string]int64 {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	result, err := p.rc.HGetAll(string(BuildPendingRecoverCacheKey(tenantId, ruleId))).Result()
	if err != nil {
		return map[string]int64{}
	}

	var newMap = make(map[string]int64)
	for k, v := range result {
		newMap[k] = tools.ConvertStringToInt64(v)
	}

	return newMap
}

func BuildPendingRecoverCacheKey(tenantId, ruleId string) PendingRecoverCacheKey {
	return PendingRecoverCacheKey(fmt.Sprintf("w8t:%s:pendingRecover:%s.fingerprints", tenantId, ruleId))
}
