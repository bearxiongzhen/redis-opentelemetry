package main

import (
	"context"
	"github.com/go-redis/redis"

	apmgoredis "github.com/bearxiongzhen/redis-opentelemetry"
)

var (
	redisClient redis.UniversalClient
)

func InitRedisConnection() {
	redisClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{"******:6380"},
		Password: "******", // no password set
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		panic(err.Error())
	}
	apmgoredis.InitTracingWrap()
}

// RedisConn 获取redis链接
func RedisConnection(ctx context.Context) redis.UniversalClient {
	return apmgoredis.Wrap(redisClient).WithContext(ctx).Cluster()
}
