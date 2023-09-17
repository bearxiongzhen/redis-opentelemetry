
# 0. 背景
提供了go-redis v7及以下版本支持opentelemetry协议的埋点功能, 这里以go-redis v6.15.5+incompatible版本为基础，

# 1. 目录结构

```
.
├── README.md
├── examples
│   └── main.go
│   └── init.go
│   └── redis_cluster_connection.go
├── go.mod
├── go.sum
├── client.go
├── rediscmd.go
├── redosptel.go
 
```

1. 示例入口见/examples/main.go

2. trace相关配置在/examples/init.go文件中

3. redis的埋点sdk设置在/examples/redis_cluster_connection.go文件中 apmgoredis.InitTracingWrap()

4. gin的埋点设置见/examples/main.go函数中

5. 每次的请求redis都会默认操作Wrap:apmgoredis.Wrap(redisClient).WithContext(ctx)，见/examples/redis_cluster_connection.go


# 2. 术语说明
**Span**：一个节点在收到请求以及完成请求的过程是一个 `Span`，`Span` 记录了在这个过程中产生的各种信息。

**Trace**：一条`Trace`（调用链）可以被认为是一个由多个`Span`组成的有向无环图（DAG图）， `Span`与`Span`的关系被命名为`References`。

**Tracer**：`Tracer`接口用来创建`Span`，以及处理如何处理`Inject`(serialize) 和 `Extract` (deserialize)，用于跨进程边界传递。

**Context**:  `Context`是一个非常重要的概念，当我们需要跨服务传播trace数据时，可以在`Context`中存储spanID，traceID等信息，并随请求传输到另一个服务。

#### step 1：配置opentelemetry上报地址、服务名、token信息，设置TracerProvide初始化

```go
//New exporter
opts := []otlptracegrpc.Option{
otlptracegrpc.WithEndpoint("127.0.0.1:4317"), // 替换成apm上报地址
otlptracegrpc.WithInsecure(),
}
exporter, err := otlptracegrpc.New(ctx, opts...)
if err != nil {
log.Fatal(err)
}

//设置Token，也可以设置环境变量：OTEL_RESOURCE_ATTRIBUTES=token=xxxxxxxxx
r, err := resource.New(ctx, []resource.Option{
resource.WithAttributes(
attribute.KeyValue{Key: "token", Value: attribute.StringValue("gFCZSIqDCUYQRAMjJSEp")},
attribute.KeyValue{Key: "service.name", Value: attribute.StringValue("Test-service")},
),
}...)
if err != nil {
log.Fatal(err)
}

//New TracerProvider
tp := sdktrace.NewTracerProvider(
sdktrace.WithSampler(sdktrace.AlwaysSample()),
sdktrace.WithBatcher(exporter),
sdktrace.WithResource(r),
)
otel.SetTracerProvider(tp)
otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
return tp
```

#### step 2：redis初始化配置，并设置opentelemtry redis相关等信息apmgoredis.InitTracingWrap()，这是redis opentelemetry协议核心初始化配置

```go
func InitRedisConnection() {
    redisClient = redis.NewClusterClient(&redis.ClusterOptions{
    Addrs:    []string{"9.135.71.56:6380"},
    Password: "1qaz2wsx", // no password set
    })
    _, err := redisClient.Ping().Result()
    if err != nil {
    panic(err.Error())
    }
	// 这是redis产生span的核心配置初始化操作，设置redis tracer配置的相关信息
    apmgoredis.InitTracingWrap()
}
```

#### step 3：产生redis span

这个时候数据上报已完成，可以在UI界面中看到了，完整代码如下

```go
// RedisConn 获取redis的链接，每次的redis请求都会组装Wrap，保证每一个redis请求都会触发产生span的操作
func RedisConnection(ctx context.Context) redis.UniversalClient {
    return apmgoredis.Wrap(redisClient).WithContext(ctx)
}
```

# 3. 相关sdk链接

1. [opentelemetry-go-instrument](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation)
2. [gorm otel](https://github.com/go-gorm/opentelemetry)

