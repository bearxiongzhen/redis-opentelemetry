module git.woa.com/taw/otel-redis/examples

go 1.18

require (
	github.com/go-redis/redis v6.15.5+incompatible
	github.com/gin-gonic/gin v1.9.1
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.32.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.7.0
    go.opentelemetry.io/otel/sdk v1.8.0
)
