package redis

import (
	"context"
	"fmt"

	rdb "github.com/redis/go-redis/v9"
)

type RedisNodesSetRepository struct {
	client *rdb.Client
}

func NewRedisNodesSetRepository(redisAddr string) *RedisNodesSetRepository {
	client := rdb.NewClient(&rdb.Options{
		Addr: redisAddr,
	})
	return &RedisNodesSetRepository{
		client: client,
	}
}

func (r *RedisNodesSetRepository) GetChatNodes(ctx context.Context, chatID string) ([]string, error) {
	nodes, err := r.client.SMembers(ctx, chatID).Result()
	if err != nil {
		return nil, fmt.Errorf("redis: GetChatNodes error: %s", err.Error())
	}
	return nodes, nil
}
