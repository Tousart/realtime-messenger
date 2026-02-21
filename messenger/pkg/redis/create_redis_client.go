package redis

import (
	"context"

	rdb "github.com/redis/go-redis/v9"
)

type RedisConnection struct {
	client *rdb.Client
	pubsub *rdb.PubSub
}

func NewClient(addr string) *RedisConnection {
	return &RedisConnection{
		client: rdb.NewClient(&rdb.Options{
			Addr: addr,
		}),
	}
}

func (c *RedisConnection) CreatePubSub(ctx context.Context) *rdb.PubSub {
	c.pubsub = c.client.Subscribe(ctx)
	return c.pubsub
}

func (c *RedisConnection) Client() *rdb.Client {
	return c.client
}

func (c *RedisConnection) PubSub() *rdb.PubSub {
	return c.pubsub
}

func (c *RedisConnection) Close() {
	if c.pubsub != nil {
		c.pubsub.Close()
	}
	if c.client != nil {
		c.client.Close()
	}
}
