package cache

import (
	"watchAlert/pkg/client"
)

type (
	entryCache struct {
		//redis    *redis.Client
		provider *ProviderPoolStore
		cache    Cache
	}

	InterEntryCache interface {
		Cache() Cache
		Silence() SilenceCacheInterface
		Event() EventCacheInterface
		ProviderPools() *ProviderPoolStore
		FaultCenter() FaultCenterCacheInterface
	}
)

func NewEntryCache(cacheType string) InterEntryCache {

	p := NewClientPoolStore()
	switch cacheType {
	case "Redis":
		return &entryCache{
			cache:    client.NewRedisCache(),
			provider: p,
		}
	default:
		return &entryCache{
			cache:    client.InitLocalCache(),
			provider: p,
		}
	}
}

func (e entryCache) Cache() Cache                      { return e.cache }
func (e entryCache) Silence() SilenceCacheInterface    { return newSilenceCacheInterface(e.cache) }
func (e entryCache) Event() EventCacheInterface        { return newEventCacheInterface(e.cache) }
func (e entryCache) ProviderPools() *ProviderPoolStore { return e.provider }
func (e entryCache) FaultCenter() FaultCenterCacheInterface {
	return newFaultCenterCacheInterface(e.cache)
}
