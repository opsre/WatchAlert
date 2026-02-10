package client

import (
	"fmt"
	"log"
	"watchAlert/config"

	"github.com/go-redis/redis"
)

var Redis *redis.Client

func InitRedis() *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Application.Redis.Host, config.Application.Redis.Port),
		Password: config.Application.Redis.Pass,
		DB:       config.Application.Redis.Database, // 使用默认的数据库
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
