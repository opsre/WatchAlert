package cache

import (
	"github.com/bytedance/sonic"
	"github.com/go-redis/redis"
	"sync"
	"time"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"
)

type (
	// ProbingCache 用于管理拨测事件缓存操作
	ProbingCache struct {
		rc *redis.Client
		sync.RWMutex
	}

	// ProbingCacheInterface 定义了事件缓存的操作接口
	ProbingCacheInterface interface {
		SetProbingEventCache(event models.ProbingEvent, expiration time.Duration)
		GetProbingEventCache(key models.ProbingEventCacheKey) (*models.ProbingEvent, error)
		DelProbingEventCache(key models.ProbingEventCacheKey) error
		GetProbingEventFirstTime(key models.ProbingEventCacheKey) int64
		GetProbingEventLastEvalTime(key models.ProbingEventCacheKey) int64
		GetProbingEventLastSendTime(key models.ProbingEventCacheKey) int64
	}
)

// newProbingCacheInterface 创建一个新的 ProbingCache 实例
func newProbingCacheInterface(r *redis.Client) ProbingCacheInterface {
	return &ProbingCache{
		rc: r,
	}
}

// SetProbingEventCache 设置探测事件缓存
func (p *ProbingCache) SetProbingEventCache(event models.ProbingEvent, expiration time.Duration) {
	eventJSON := tools.JsonMarshalToString(event)
	p.setProbingCache(models.BuildProbingEventCacheKey(event.TenantId, event.RuleId), eventJSON, expiration)
}

// GetProbingEventCache 获取探测事件缓存
func (p *ProbingCache) GetProbingEventCache(key models.ProbingEventCacheKey) (*models.ProbingEvent, error) {
	var event *models.ProbingEvent

	data, err := p.getProbingCache(key)
	if err != nil {
		return event, err
	}

	if err := sonic.Unmarshal([]byte(data), &event); err != nil {
		return event, err
	}

	return event, nil
}

func (p *ProbingCache) DelProbingEventCache(key models.ProbingEventCacheKey) error {
	return p.delProbingCache(key)
}

// GetProbingEventFirstTime 获取探测事件的首次触发时间
func (p *ProbingCache) GetProbingEventFirstTime(key models.ProbingEventCacheKey) int64 {
	event, err := p.GetProbingEventCache(key)
	if err != nil || event.FirstTriggerTime == 0 {
		return time.Now().Unix()
	}
	return event.FirstTriggerTime
}

// GetProbingEventLastEvalTime 获取探测事件的最后评估时间
func (p *ProbingCache) GetProbingEventLastEvalTime(key models.ProbingEventCacheKey) int64 {
	curTime := time.Now().Unix()
	event, err := p.GetProbingEventCache(key)
	if err != nil || event.LastEvalTime == 0 || event.LastEvalTime < curTime {
		return curTime
	}
	return event.LastEvalTime
}

// GetProbingEventLastSendTime 获取探测事件的最后发送时间
func (p *ProbingCache) GetProbingEventLastSendTime(key models.ProbingEventCacheKey) int64 {
	event, err := p.GetProbingEventCache(key)
	if err != nil {
		return 0
	}
	return event.LastSendTime
}

// 封装 Redis 操作
func (p *ProbingCache) setProbingCache(key models.ProbingEventCacheKey, value string, expiration time.Duration) {
	p.rc.Set(string(key), value, expiration)
}

func (p *ProbingCache) getProbingCache(key models.ProbingEventCacheKey) (string, error) {
	return p.rc.Get(string(key)).Result()
}

func (p *ProbingCache) delProbingCache(key models.ProbingEventCacheKey) error {
	return p.rc.Del(string(key)).Err()
}
