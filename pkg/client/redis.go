package client

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
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
