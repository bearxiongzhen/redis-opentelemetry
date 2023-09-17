package apm

import (
	"context"
	
	"github.com/go-redis/redis"
)

type Client interface {
	redis.UniversalClient

	// ClusterClient returns the wrapped *redis.ClusterClient,;
	// or nil if a non-cluster client is wrapped.;
	Cluster() *redis.ClusterClient

	// Ring returns the wrapped *redis.Ring,;
	// or nil if a non-ring client is wrapped.;
	RingClient() *redis.Ring

	// WithContext returns a shallow copy of the client with
	// its context changed to ctx and will add instrumentation
	// with client.WrapProcess and client.WrapProcessPipeline
	//
	// To report commands as spans, ctx must contain a transaction or span.
	WithContext(ctx context.Context) Client
}

// Wrap wraps client such that executed commands are reported as spans to Elastic APM,
// using the client's associated context.
// A context-specific client may be obtained by using Client.WithContext.
func Wrap(client redis.UniversalClient) Client {
	switch client.(type) {
	case *redis.Client:
		return contextClient{Client: client.(*redis.Client)}
	case *redis.ClusterClient:
		return contextClusterClient{ClusterClient: client.(*redis.ClusterClient)}
	case *redis.Ring:
		return contextRingClient{Ring: client.(*redis.Ring)}
	}

	return client.(Client)

}

type contextClient struct {
	*redis.Client
}

func (c contextClient) WithContext(ctx context.Context) Client {
	c.Client = c.Client.WithContext(ctx)

	c.WrapProcess(process(ctx))
	c.WrapProcessPipeline(processPipeline(ctx))

	return c
}

func (c contextClient) Cluster() *redis.ClusterClient {
	return nil
}

func (c contextClient) RingClient() *redis.Ring {
	return nil
}

type contextClusterClient struct {
	*redis.ClusterClient
}

func (c contextClusterClient) Cluster() *redis.ClusterClient {
	return c.ClusterClient
}

func (c contextClusterClient) RingClient() *redis.Ring {
	return nil
}

func (c contextClusterClient) WithContext(ctx context.Context) Client {
	c.ClusterClient = c.ClusterClient.WithContext(ctx)

	c.WrapProcess(process(ctx))
	c.WrapProcessPipeline(processPipeline(ctx))

	return c
}

type contextRingClient struct {
	*redis.Ring
}

func (c contextRingClient) Cluster() *redis.ClusterClient {

	return nil
}

func (c contextRingClient) RingClient() *redis.Ring {
	return c.Ring
}

func (c contextRingClient) WithContext(ctx context.Context) Client {
	c.Ring = c.Ring.WithContext(ctx)

	c.WrapProcess(process(ctx))
	c.WrapProcessPipeline(processPipeline(ctx))

	return c
}

func process(ctx context.Context) func(oldProcess func(cmd redis.Cmder) error) func(cmd redis.Cmder) error {
	return func(oldProcess func(cmd redis.Cmder) error) func(cmd redis.Cmder) error {
		return func(cmd redis.Cmder) error {
			if TracingWrap != nil {
				ctx, _ = TracingWrap.BeforeProcess(ctx, cmd)
				defer TracingWrap.AfterProcess(ctx, cmd)
			}
			return oldProcess(cmd)
		}
	}
}

func processPipeline(ctx context.Context) func(oldProcess func(cmds []redis.Cmder) error) func(cmds []redis.Cmder) error {
	return func(oldProcess func(cmds []redis.Cmder) error) func(cmds []redis.Cmder) error {
		return func(cmds []redis.Cmder) error {
			if TracingWrap != nil {
				ctx, _ = TracingWrap.BeforeProcessPipeline(ctx, cmds)
				defer TracingWrap.AfterProcessPipeline(ctx, cmds)
			}

			return oldProcess(cmds)
		}
	}
}
