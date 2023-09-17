package main

import (
	"context"
	"github.com/go-redis/redis"

	apmgoredis "github.com/opentelemetry/goredis"
)

var (
	redisClient redis.UniversalClient
)

func InitRedisConnection() {
	redisClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{"9.135.71.56:6380"},
		Password: "1qaz2wsx", // no password set
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		panic(err.Error())
	}
	apmgoredis.InitTracingWrap()
}

// RedisConn 获取redis链接
func RedisConnection(ctx context.Context) redis.UniversalClient {
	return apmgoredis.Wrap(redisClient).WithContext(ctx)
}
