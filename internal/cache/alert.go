package cache

import (
	"context"
	"fmt"
	"sync"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"

	"github.com/bytedance/sonic"
	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/core/logc"
)

type (
	// AlertCache 用于管理告警事件缓存操作
	AlertCache struct {
		rc *redis.Client
		sync.RWMutex
	}

	// AlertCacheInterface 定义了事件缓存的操作接口
	AlertCacheInterface interface {
		PushAlertEvent(event *models.AlertCurEvent)
		RemoveAlertEvent(tenantId, faultCenterId, fingerprint string)
		GetFingerprintsByRuleId(tenantId, faultCenterId, ruleId string) []string
		GetAllEvents(key models.AlertEventCacheKey) (map[string]*models.AlertCurEvent, error)
		GetEventFromCache(tenantId, faultCenterId, fingerprint string) (models.AlertCurEvent, error)
	}
)

// newAlertCacheInterface 创建一个新的 AlertCache 实例
func newAlertCacheInterface(r *redis.Client) AlertCacheInterface {
	return &AlertCache{
		rc: r,
	}
}

// PushAlertEvent 将事件推送到故障中心的缓存中
func (a *AlertCache) PushAlertEvent(event *models.AlertCurEvent) {
	key := models.BuildAlertEventCacheKey(event.TenantId, event.FaultCenterId)
	a.setEventCacheHash(key, event.Fingerprint, tools.JsonMarshalToString(event))
}

// RemoveAlertEvent 从故障中心的缓存中移除事件
func (a *AlertCache) RemoveAlertEvent(tenantId, faultCenterId, fingerprint string) {
	key := models.BuildAlertEventCacheKey(tenantId, faultCenterId)
	a.deleteEventCacheHash(key, fingerprint)
}

// GetAllEvents 获取故障中心的所有事件
func (a *AlertCache) GetAllEvents(key models.AlertEventCacheKey) (map[string]*models.AlertCurEvent, error) {
	a.RLock()
	defer a.RUnlock()

	result, err := a.getEventCacheHashAll(key)
	if err != nil {
		return nil, err
	}

	events := make(map[string]*models.AlertCurEvent)
	for fingerprint, eventJSON := range result {
		var event models.AlertCurEvent
		if err := sonic.Unmarshal([]byte(eventJSON), &event); err != nil {
			logc.Error(context.Background(), fmt.Sprintf("unmarshal event json error: %s, event json: %s", err.Error(), eventJSON))
			continue
		}
		events[fingerprint] = &event
	}

	return events, nil
}

// GetFingerprintsByRuleId 获取与指定规则 ID 相关的指纹列表
func (a *AlertCache) GetFingerprintsByRuleId(tenantId, faultCenterId, ruleId string) []string {
	key := models.BuildAlertEventCacheKey(tenantId, faultCenterId)
	events, err := a.GetAllEvents(key)
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
func (a *AlertCache) GetEventFromCache(tenantId, faultCenterId, fingerprint string) (models.AlertCurEvent, error) {
	key := models.BuildAlertEventCacheKey(tenantId, faultCenterId)
	data, err := a.getEventCacheHash(key, fingerprint)
	if err != nil {
		return models.AlertCurEvent{}, err
	}

	var event models.AlertCurEvent
	if err := sonic.Unmarshal([]byte(data), &event); err != nil {
		return models.AlertCurEvent{}, err
	}

	return event, nil
}

// 封装 Redis 操作
func (a *AlertCache) setEventCacheHash(key models.AlertEventCacheKey, field, value string) {
	a.rc.HSet(string(key), field, value)
}

func (a *AlertCache) deleteEventCacheHash(key models.AlertEventCacheKey, field string) {
	a.rc.HDel(string(key), field)
}

func (a *AlertCache) getEventCacheHash(key models.AlertEventCacheKey, field string) (string, error) {
	return a.rc.HGet(string(key), field).Result()
}

func (a *AlertCache) getEventCacheHashAll(key models.AlertEventCacheKey) (map[string]string, error) {
	return a.rc.HGetAll(string(key)).Result()
}
