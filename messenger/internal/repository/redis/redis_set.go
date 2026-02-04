package redis

import (
	"context"
	"fmt"

	rdb "github.com/redis/go-redis/v9"
)

type RedisNodesSetRepository struct {
	client   *rdb.Client
	nodeAddr string
}

func NewRedisNodesSetRepository(redisAddr, nodeAddr string) *RedisNodesSetRepository {
	client := rdb.NewClient(&rdb.Options{
		Addr: redisAddr,
	})
	return &RedisNodesSetRepository{
		client:   client,
		nodeAddr: nodeAddr,
	}
}

func (r *RedisNodesSetRepository) AddNodeToChat(ctx context.Context, chatID string) error {
	err := r.client.SAdd(ctx, chatID, r.nodeAddr).Err()
	if err != nil {
		return fmt.Errorf("redis: AddNodeToChat error: %s", err.Error())
	}
	return nil
}

func (r *RedisNodesSetRepository) RemoveNodeFromChat(ctx context.Context, chatID string) error {
	err := r.client.SRem(ctx, chatID, r.nodeAddr).Err()
	if err != nil {
		return fmt.Errorf("redis: RemoveNodeFromChat error: %s", err.Error())
	}
	return nil
}
