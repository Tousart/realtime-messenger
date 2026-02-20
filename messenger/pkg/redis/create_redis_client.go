package redis

import (
	"context"

	rdb "github.com/redis/go-redis/v9"
)

func CreateRedisClient(addr string) *rdb.Client {
	return rdb.NewClient(&rdb.Options{
		Addr: addr,
	})
}

func CreateRedisPubSubObject(ctx context.Context, client *rdb.Client) *rdb.PubSub {
	return client.Subscribe(ctx)
}
