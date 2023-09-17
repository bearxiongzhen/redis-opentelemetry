package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"log"
)

func main() {

	InitTraceProvider(context.Background())

	InitRedisConnection()

	r := gin.New()

	r.Use(otelgin.Middleware("redis-16.15.5-gin"))

	r.GET("/ping", func(ginCtx *gin.Context) {
		redisClient = RedisConnection(ginCtx.Request.Context())
		redisClient.Ping()
	})

	r.GET("/set", func(ginCtx *gin.Context) {
		redisClient = RedisConnection(ginCtx.Request.Context())
		redisClient.Set("data", "world", 1000000) //1ms
		data := redisClient.Get("data").String()
		log.Println(data)
	})

	r.GET("/pipeline", func(ginCtx *gin.Context) {
		redisClient = RedisConnection(ginCtx.Request.Context())
		pipeline := redisClient.Pipeline()

		for index := 0; index < 10; index++ {
			pipeline.Get("1")
		}
		pipeline.Exec()
	})
	_ = r.Run(":8082")
}
