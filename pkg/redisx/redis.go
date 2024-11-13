package redisx

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"os"
	"sync"
	"time"
	"wallet-service/pkg/config"
	"wallet-service/pkg/logger"
)

var one = sync.Once{}
var client *redis.Client

func GetRedisClient() *redis.Client {
	one.Do(func() {
		if client == nil {
			newClient()
		}
	})

	return client
}

func InitRedis() {
	one.Do(func() {
		if client == nil {
			newClient()
		}
	})
}

func newClient() *redis.Client {
	var err error
	var ctx = context.Background()
	conf := config.GetConfig()

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost != "" {
		conf.Redis.Host = redisHost
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort != "" {
		conf.Redis.Port = cast.ToInt(redisPort)
	}

	client = redis.NewClient(&redis.Options{
		PoolSize:        conf.Redis.PoolSize,                                     // 连接池的大小
		PoolTimeout:     time.Duration(conf.Redis.PoolTimeout) * time.Second,     // 连接池内获取可用连接超时
		MinIdleConns:    conf.Redis.MinIdleConns,                                 // 最小空闲连接数
		MaxIdleConns:    conf.Redis.MaxIdleConns,                                 // 最大空闲连接数
		ConnMaxIdleTime: time.Duration(conf.Redis.ConnMaxIdleTime) * time.Second, // 连接的最大空闲时间
		Addr:            fmt.Sprintf("%s:%d", conf.Redis.Host, conf.Redis.Port),
		Password:        conf.Redis.Password,
		DB:              conf.Redis.DB,
	})

	_, err = client.Ping(ctx).Result()
	if err != nil {
		logger.Panic(ctx, "redis Init err : %v", err)
		panic(err)
	}
	return client
}
