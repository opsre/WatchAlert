package client

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"time"
	"watchAlert/internal/global"
)

var Redis *redis.Client

func InitRedis() *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", global.Config.Redis.Host, global.Config.Redis.Port),
		Password: global.Config.Redis.Pass,
		DB:       global.Config.Redis.Database, // 使用默认的数据库
	})

	// 尝试连接到 Redis 服务器
	_, err := client.Ping().Result()
	if err != nil {
		log.Printf("redis Connection Failed %s", err)
		panic(err)
	}

	Redis = client

	return client

}

type RedisCache struct {
	*redis.Client
}

func NewRedisCache() *RedisCache {
	return &RedisCache{InitRedis()}
}

func (c *RedisCache) SetKey(key, value string, expiration time.Duration) {
	c.Set(key, value, expiration)
}

func (c *RedisCache) GetKey(key string) (string, error) {
	return c.Get(key).Result()
}

func (c *RedisCache) SetHashAny(key, field string, value any) {
	c.HSet(key, field, value)
}

func (c *RedisCache) DeleteKey(key string) {
	c.Del(key)
}

func (c *RedisCache) SetHash(key, field, value string) {
	c.HSet(key, field, value)
}

func (c *RedisCache) DeleteHash(key, field string) {
	c.HDel(key, field)
}

func (c *RedisCache) GetHash(key, field string) (string, error) {
	return c.HGet(key, field).Result()
}

func (c *RedisCache) GetHashAll(key string) (map[string]string, error) {
	return c.HGetAll(key).Result()
}
