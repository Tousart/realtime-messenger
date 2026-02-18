package redis

import (
	"context"
	"fmt"
	"strconv"

	rdb "github.com/redis/go-redis/v9"
	"github.com/tousart/messenger/internal/domain"
)

type RedisQueuesRepository struct {
	client    *rdb.Client
	queueName string
}

func NewRedisQueuesRepository(client *rdb.Client, queueName string) *RedisQueuesRepository {
	return &RedisQueuesRepository{
		client:    client,
		queueName: queueName,
	}
}

func (r *RedisQueuesRepository) Queues(ctx context.Context, chat *domain.Chat) ([]string, error) {
	queues := r.client.SMembers(ctx, strconv.Itoa(chat.ChatID))
	if queues.Err() != nil {
		return nil, fmt.Errorf("redis: Queues error: %w", queues.Err())
	}
	return queues.Val(), nil
}

func (r *RedisQueuesRepository) AddQueueToChat(ctx context.Context, chat *domain.Chat) error {
	err := r.client.SAdd(ctx, strconv.Itoa(chat.ChatID), r.queueName).Err()
	if err != nil {
		return fmt.Errorf("redis: AddQueueToChat error: %w", err)
	}
	return nil
}

func (r *RedisQueuesRepository) RemoveQueueFromChat(ctx context.Context, chat *domain.Chat) error {
	err := r.client.SRem(ctx, strconv.Itoa(chat.ChatID), r.queueName).Err()
	if err != nil {
		return fmt.Errorf("redis: RemoveQueueFromChat error: %w", err)
	}
	return nil
}
