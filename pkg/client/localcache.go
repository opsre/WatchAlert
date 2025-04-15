package client

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"github.com/zeromicro/go-zero/core/logc"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	LocalCache *Cache
	once       sync.Once
	cacheFile  = "cache_file"
)

type Cache struct {
	cache    *cache.Cache
	mu       sync.Mutex
	shutdown chan struct{}
	//wg       sync.WaitGroup
}

func InitLocalCache() *Cache {
	var localCache *Cache
	once.Do(
		func() {
			c := cache.New(cache.DefaultExpiration, cache.NoExpiration)
			gob.Register(cache.Item{})
			gob.Register(map[string]interface{}{})
			gob.Register(map[string]string{})

			if err := c.LoadFile(cacheFile); err != nil {
				logc.Debugf(context.Background(), "缓存加载失败，首次启动请忽略。%v", err)
			} else {
				logc.Debugf(context.Background(), "加载缓存成功")
			}

			localCache = &Cache{cache: c, shutdown: make(chan struct{})}
			LocalCache = localCache
		})

	//localCache.wg.Add(1)
	go localCache.persist()

	//localCache.wg.Add(1)
	go localCache.listenSignals()

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

func (c *Cache) listenSignals() {
	//defer c.wg.Done()

	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	select {
	case <-sigChan:
		logc.Infof(context.Background(), "接收到终止信号，启动关闭流程")
		close(c.shutdown)
	case <-c.shutdown:
	}
}

func (c *Cache) persist() {
	//defer c.wg.Done()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := c.cache.SaveFile(cacheFile); err != nil {
				log.Warn(err)
			}
			logc.Debugf(context.Background(), "定时缓存持久化成功")
		case <-c.shutdown:
			if err := c.cache.SaveFile(cacheFile); err != nil {
				logc.Errorf(context.Background(), fmt.Sprintf("退出前缓存持久化异常: %v", err))
			}
			logc.Debugf(context.Background(), "缓存持久化成功")
			return
		}
	}
}
