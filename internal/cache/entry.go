package cache

import (
	"github.com/go-redis/redis"
	"watchAlert/pkg/client"
)

type (
	entryCache struct {
		redis    *redis.Client
		provider *ProviderPoolStore
	}

	InterEntryCache interface {
		Redis() *redis.Client
		Silence() SilenceCacheInterface
		Event() EventCacheInterface
		ProviderPools() *ProviderPoolStore
		FaultCenter() FaultCenterCacheInterface
	}
)

func NewEntryCache() InterEntryCache {
	r := client.InitRedis()
	p := NewClientPoolStore()

	return &entryCache{
		redis:    r,
		provider: p,
	}
}

func (e entryCache) Redis() *redis.Client              { return e.redis }
func (e entryCache) Silence() SilenceCacheInterface    { return newSilenceCacheInterface(e.redis) }
func (e entryCache) Event() EventCacheInterface        { return newEventCacheInterface(e.redis) }
func (e entryCache) ProviderPools() *ProviderPoolStore { return e.provider }
func (e entryCache) FaultCenter() FaultCenterCacheInterface {
	return newFaultCenterCacheInterface(e.redis)
}
