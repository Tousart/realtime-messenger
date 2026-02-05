package redis

import (
	"context"
	"fmt"

	rdb "github.com/redis/go-redis/v9"
)

type RedisNodesReceiverRepository struct {
	client   *rdb.Client
	nodeAddr string
}

func NewRedisNodesReceiverRepository(redisAddr, nodeAddr string) *RedisNodesReceiverRepository {
	client := rdb.NewClient(&rdb.Options{
		Addr: redisAddr,
	})
	return &RedisNodesReceiverRepository{
		client:   client,
		nodeAddr: nodeAddr,
	}
}

func (r *RedisNodesReceiverRepository) AddNodeToChat(ctx context.Context, chatID string) error {
	err := r.client.SAdd(ctx, chatID, r.nodeAddr).Err()
	if err != nil {
		return fmt.Errorf("redis: AddNodeToChat error: %s", err.Error())
	}
	return nil
}

func (r *RedisNodesReceiverRepository) RemoveNodeFromChat(ctx context.Context, chatID string) error {
	err := r.client.SRem(ctx, chatID, r.nodeAddr).Err()
	if err != nil {
		return fmt.Errorf("redis: RemoveNodeFromChat error: %s", err.Error())
	}
	return nil
}
