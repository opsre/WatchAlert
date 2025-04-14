package client

import (
	"errors"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

var LocalCache *Cache
var once sync.Once

type Cache struct {
	cache *cache.Cache
	mu    sync.Mutex
}

var (
	ErrKeyNotFound   = errors.New("key not found")
	ErrFieldNotFound = errors.New("field not found")
	ErrInvalidType   = errors.New("invalid type")
)

func InitLocalCache() *Cache {
	var localCache *Cache
	once.Do(
		func() {
			c := cache.New(cache.DefaultExpiration, cache.NoExpiration)
			localCache = &Cache{cache: c}
			LocalCache = localCache
		})
	return localCache
}

func (c *Cache) SetKey(key, value string, expiration time.Duration) {
	c.cache.Set(key, value, expiration)
}

func (c *Cache) GetKey(key string) (string, error) {
	value, exist := c.cache.Get(key)
	if !exist {
		return "", nil
	}
	strValue, ok := value.(string)
	if !ok {
		return "", nil
	}
	return strValue, nil
}

func (c *Cache) DeleteKey(key string) {
	c.cache.Delete(key)
}

func (c *Cache) SetHash(key, field, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var hashMap map[string]string
	if val, found := c.cache.Get(key); found {
		hashMap, _ = val.(map[string]string)
	}
	if hashMap == nil {
		hashMap = make(map[string]string)
	}

	hashMap[field] = value
	c.cache.Set(key, hashMap, cache.NoExpiration)
}

func (c *Cache) SetHashAny(key, field string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var hashMap map[string]interface{}
	if val, found := c.cache.Get(key); found {
		hashMap, _ = val.(map[string]interface{})
	}
	if hashMap == nil {
		hashMap = make(map[string]interface{})
	}

	hashMap[field] = value
	c.cache.Set(key, hashMap, cache.NoExpiration)
}

func (c *Cache) DeleteHash(key, field string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if val, found := c.cache.Get(key); found {
		if hashMap, ok := val.(map[string]string); ok {
			delete(hashMap, field)
			c.cache.Set(key, hashMap, cache.NoExpiration)
		}
	}
}

func (c *Cache) GetHash(key, field string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, found := c.cache.Get(key)
	if !found {
		return "", nil
	}

	hashMap, ok := val.(map[string]string)
	if !ok {
		return "", nil
	}

	value, exist := hashMap[field]
	if !exist {
		return "", nil
	}

	return value, nil
}

func (c *Cache) GetHashAll(key string) (map[string]string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, found := c.cache.Get(key)
	if !found {
		return nil, nil
	}

	hashMap, ok := val.(map[string]string)
	if !ok {
		return nil, nil
	}

	return hashMap, nil
}
