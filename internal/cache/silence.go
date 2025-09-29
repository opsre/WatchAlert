package cache

import (
	"github.com/bytedance/sonic"
	"github.com/go-redis/redis"
	"sync"
	"watchAlert/internal/models"
	"watchAlert/pkg/client"
	"watchAlert/pkg/tools"
)

type (
	// SilenceCache 用于管理告警静默的缓存操作
	SilenceCache struct {
		rc *redis.Client
		sync.RWMutex
	}

	// SilenceCacheInterface 定义了告警静默缓存的操作接口
	SilenceCacheInterface interface {
		PushAlertMute(mute models.AlertSilences)
		RemoveAlertMute(tenantId, faultCenterId, id string)
		GetAlertMutes(tenantId, faultCenterId string) ([]string, error)
		WithIdGetMuteFromCache(tenantId, faultCenterId, id string) (*models.AlertSilences, error)
	}
)

// newSilenceCacheInterface 创建一个新的 SilenceCache 实例
func newSilenceCacheInterface(r *redis.Client) SilenceCacheInterface {
	return &SilenceCache{
		rc: r,
	}
}

// PushAlertMute 将静默规则推送到故障中心的缓存中
func (sc *SilenceCache) PushAlertMute(mute models.AlertSilences) {
	sc.Lock()
	defer sc.Unlock()

	key := models.BuildAlertMuteCacheKey(mute.TenantId, mute.FaultCenterId)
	sc.setRedisHash(key, mute.ID, tools.JsonMarshalToString(mute))
}

// RemoveAlertMute 从故障中心的缓存中移除静默规则
func (sc *SilenceCache) RemoveAlertMute(tenantId, faultCenterId, id string) {
	sc.Lock()
	defer sc.Unlock()

	key := models.BuildAlertMuteCacheKey(tenantId, faultCenterId)
	sc.deleteRedisHash(key, id)
}

func (sc *SilenceCache) GetAlertMutes(tenantId, faultCenterId string) ([]string, error) {
	sc.RLock()
	defer sc.RUnlock()

	key := models.BuildAlertMuteCacheKey(tenantId, faultCenterId)
	mapping, err := sc.getRedisAllHashMap(key)
	if err != nil {
		return nil, err
	}
	var ids []string
	for id := range mapping {
		ids = append(ids, id)
	}
	return ids, nil
}

// WithIdGetMuteFromCache 从缓存中获取静默规则
func (sc *SilenceCache) WithIdGetMuteFromCache(tenantId, faultCenterId, id string) (*models.AlertSilences, error) {
	key := models.BuildAlertMuteCacheKey(tenantId, faultCenterId)
	cache, err := sc.getRedisHash(key, id)
	if err != nil {
		return nil, err
	}

	var mute models.AlertSilences
	if err := sonic.Unmarshal(cache, &mute); err != nil {
		return nil, err
	}

	return &mute, nil
}

// setRedisHash 设置 Redis 哈希表中的值
func (sc *SilenceCache) setRedisHash(key models.AlertMuteCacheKey, field string, value interface{}) {
	client.Redis.HSet(string(key), field, value)
}

// deleteRedisHash 删除 Redis 哈希表中的值
func (sc *SilenceCache) deleteRedisHash(key models.AlertMuteCacheKey, field string) {
	client.Redis.HDel(string(key), field)
}

// getRedisHash 获取 Redis 哈希表中的值
func (sc *SilenceCache) getRedisHash(key models.AlertMuteCacheKey, field string) ([]byte, error) {
	return sc.rc.HGet(string(key), field).Bytes()
}

// getRedisAllMap 获取 Redis 哈希表Map
func (sc *SilenceCache) getRedisAllHashMap(key models.AlertMuteCacheKey) (map[string]string, error) {
	return sc.rc.HGetAll(string(key)).Result()
}
