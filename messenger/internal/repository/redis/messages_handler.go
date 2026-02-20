package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisMessagesHandlerRepository struct {
	client *redis.Client
	pubsub *redis.PubSub
}

func NewRedisMessagesHandlerRepository(client *redis.Client, pubsub *redis.PubSub) *RedisMessagesHandlerRepository {
	return &RedisMessagesHandlerRepository{
		client: client,
		pubsub: pubsub,
	}
}

func (mh *RedisMessagesHandlerRepository) PublishMessageToChat(ctx context.Context, chatID string, messagePayload []byte) error {
	log.Printf("сообщение опубликовано в repository\n")

	if err := mh.client.Publish(ctx, chatID, messagePayload).Err(); err != nil {
		return fmt.Errorf("redis: PublishMessageToChat: %w", err)
	}
	return nil
}

func (mh *RedisMessagesHandlerRepository) SubscribeToChats(ctx context.Context, chatIDs ...string) error {
	if err := mh.pubsub.Subscribe(ctx, chatIDs...); err != nil {
		return fmt.Errorf("redis: SubscribeToChats: %w", err)
	}
	return nil
}

func (mh *RedisMessagesHandlerRepository) UnsubscribeFromChats(ctx context.Context, chatIDs ...string) error {
	if err := mh.pubsub.Unsubscribe(ctx, chatIDs...); err != nil {
		return fmt.Errorf("redis: UnsubscribeFromChats: %w", err)
	}
	return nil
}
