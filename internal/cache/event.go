package cache

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/core/logc"
	"sync"
	"time"
	"watchAlert/internal/models"
	"watchAlert/pkg/client"
	"watchAlert/pkg/tools"
)

type (
	// EventCache 用于管理事件缓存操作
	EventCache struct {
		rc *redis.Client
		sync.RWMutex
	}

	// EventCacheInterface 定义了事件缓存的操作接口
	EventCacheInterface interface {
		SetProbingEventCache(event models.ProbingEvent, expiration time.Duration)
		GetProbingEventCache(key string) (models.ProbingEvent, error)
		GetProbingEventFirstTime(key string) int64
		GetProbingEventLastEvalTime(key string) int64
		GetProbingEventLastSendTime(key string) int64
		PushEventToFaultCenter(event *models.AlertCurEvent)
		RemoveEventFromFaultCenter(tenantId, faultCenterId, fingerprint string)
		GetFingerprintsByRuleId(tenantId, faultCenterId, ruleId string) []string
		GetAllEventsForFaultCenter(fcKey string) (map[string]models.AlertCurEvent, error)
		GetFirstTimeForFaultCenter(tenantId, faultCenterId, fingerprint string) int64
		GetLastEvalTimeForFaultCenter() int64
		GetLastSendTimeForFaultCenter(tenantId, faultCenterId, fingerprint string) int64
		GetEventStatusForFaultCenter(tenantId, faultCenterId, fingerprint string) int64
	}
)

// newEventCacheInterface 创建一个新的 EventCache 实例
func newEventCacheInterface(r *redis.Client) EventCacheInterface {
	return &EventCache{
		rc: r,
	}
}

// SetProbingEventCache 设置探测事件缓存
func (ec *EventCache) SetProbingEventCache(event models.ProbingEvent, expiration time.Duration) {
	ec.Lock()
	defer ec.Unlock()

	eventJSON, _ := json.Marshal(event)
	ec.setRedisKey(event.GetFiringAlertCacheKey(), string(eventJSON), expiration)
}

// GetProbingEventCache 获取探测事件缓存
func (ec *EventCache) GetProbingEventCache(key string) (models.ProbingEvent, error) {
	var event models.ProbingEvent

	data, err := ec.getRedisKey(key)
	if err != nil {
		return event, err
	}

	if err := json.Unmarshal([]byte(data), &event); err != nil {
		return event, err
	}

	return event, nil
}

// GetProbingEventFirstTime 获取探测事件的首次触发时间
func (ec *EventCache) GetProbingEventFirstTime(key string) int64 {
	event, err := ec.GetProbingEventCache(key)
	if err != nil || event.FirstTriggerTime == 0 {
		return time.Now().Unix()
	}
	return event.FirstTriggerTime
}

// GetProbingEventLastEvalTime 获取探测事件的最后评估时间
func (ec *EventCache) GetProbingEventLastEvalTime(key string) int64 {
	curTime := time.Now().Unix()
	event, err := ec.GetProbingEventCache(key)
	if err != nil || event.LastEvalTime == 0 || event.LastEvalTime < curTime {
		return curTime
	}
	return event.LastEvalTime
}

// GetProbingEventLastSendTime 获取探测事件的最后发送时间
func (ec *EventCache) GetProbingEventLastSendTime(key string) int64 {
	event, err := ec.GetProbingEventCache(key)
	if err != nil {
		return 0
	}
	return event.LastSendTime
}

// PushEventToFaultCenter 将事件推送到故障中心的缓存中
func (ec *EventCache) PushEventToFaultCenter(event *models.AlertCurEvent) {
	ec.Lock()
	defer ec.Unlock()

	key := models.BuildCacheEventKey(event.TenantId, event.FaultCenterId)
	ec.setRedisHash(key, event.Fingerprint, tools.JsonMarshal(event))
}

// RemoveEventFromFaultCenter 从故障中心的缓存中移除事件
func (ec *EventCache) RemoveEventFromFaultCenter(tenantId, faultCenterId, fingerprint string) {
	ec.Lock()
	defer ec.Unlock()

	key := models.BuildCacheEventKey(tenantId, faultCenterId)
	ec.deleteRedisHash(key, fingerprint)
}

// GetAllEventsForFaultCenter 获取故障中心的所有事件
func (ec *EventCache) GetAllEventsForFaultCenter(fcKey string) (map[string]models.AlertCurEvent, error) {
	ec.RLock()
	defer ec.RUnlock()

	result, err := ec.getRedisHashAll(fcKey)
	if err != nil {
		return nil, err
	}

	events := make(map[string]models.AlertCurEvent)
	for fingerprint, eventJSON := range result {
		var event models.AlertCurEvent
		if err := json.Unmarshal([]byte(eventJSON), &event); err != nil {
			return nil, err
		}
		events[fingerprint] = event
	}

	return events, nil
}

// GetFingerprintsByRuleId 获取与指定规则 ID 相关的指纹列表
func (ec *EventCache) GetFingerprintsByRuleId(tenantId, faultCenterId, ruleId string) []string {
	key := models.BuildCacheEventKey(tenantId, faultCenterId)
	events, err := ec.GetAllEventsForFaultCenter(key)
	if err != nil {
		logc.Error(context.Background(), err.Error())
		return nil
	}

	var fingerprints []string
	for fingerprint, event := range events {
		if event.RuleId == ruleId {
			fingerprints = append(fingerprints, fingerprint)
		}
	}
	return fingerprints
}

// GetEventFromCache 从缓存中获取事件数据
func (ec *EventCache) GetEventFromCache(tenantId, faultCenterId, fingerprint string) (models.AlertCurEvent, error) {
	key := models.BuildCacheEventKey(tenantId, faultCenterId)
	data, err := ec.getRedisHash(key, fingerprint)
	if err != nil {
		return models.AlertCurEvent{}, err
	}

	var event models.AlertCurEvent
	if err := json.Unmarshal([]byte(data), &event); err != nil {
		return models.AlertCurEvent{}, err
	}

	return event, nil
}

// GetFirstTimeForFaultCenter 获取故障中心事件的首次触发时间
func (ec *EventCache) GetFirstTimeForFaultCenter(tenantId, faultCenterId, fingerprint string) int64 {
	event, err := ec.GetEventFromCache(tenantId, faultCenterId, fingerprint)
	if err != nil || event.FirstTriggerTime == 0 {
		return time.Now().Unix()
	}
	return event.FirstTriggerTime
}

// GetLastEvalTimeForFaultCenter 获取故障中心事件的最后评估时间
func (ec *EventCache) GetLastEvalTimeForFaultCenter() int64 {
	return time.Now().Unix()
}

// GetLastSendTimeForFaultCenter 获取故障中心事件的最后发送时间
func (ec *EventCache) GetLastSendTimeForFaultCenter(tenantId, faultCenterId, fingerprint string) int64 {
	event, err := ec.GetEventFromCache(tenantId, faultCenterId, fingerprint)
	if err != nil {
		return 0
	}
	return event.LastSendTime
}

// GetEventStatusForFaultCenter 获取事件状态
func (ec *EventCache) GetEventStatusForFaultCenter(tenantId, faultCenterId, fingerprint string) int64 {
	event, err := ec.GetEventFromCache(tenantId, faultCenterId, fingerprint)
	if err != nil {
		return 0
	}
	return event.Status
}

// 封装 Redis 操作
func (ec *EventCache) setRedisKey(key, value string, expiration time.Duration) {
	ec.rc.Set(key, value, expiration)
}

func (ec *EventCache) getRedisKey(key string) (string, error) {
	return ec.rc.Get(key).Result()
}

func (ec *EventCache) setRedisHash(key, field, value string) {
	client.Redis.HSet(key, field, value)
}

func (ec *EventCache) deleteRedisHash(key, field string) {
	client.Redis.HDel(key, field)
}

func (ec *EventCache) getRedisHash(key, field string) (string, error) {
	return ec.rc.HGet(key, field).Result()
}

func (ec *EventCache) getRedisHashAll(key string) (map[string]string, error) {
	return ec.rc.HGetAll(key).Result()
}
