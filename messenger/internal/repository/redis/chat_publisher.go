package redis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type ChatPublisher struct {
	client *redis.Client
	pubsub *redis.PubSub
}

func NewChatPublisher(client *redis.Client, pubsub *redis.PubSub) *ChatPublisher {
	return &ChatPublisher{
		client: client,
		pubsub: pubsub,
	}
}

func (mh *ChatPublisher) PublishMessage(ctx context.Context, chatID int64, msgBytes []byte) error {
	if err := mh.client.Publish(ctx, idToString(chatID), msgBytes).Err(); err != nil {
		return fmt.Errorf("repository: redis: PublishMessage: %w", err)
	}
	return nil
}

func (mh *ChatPublisher) Subscribe(ctx context.Context, chatIDs ...int64) error {
	if chatIDs == nil {
		return nil
	}

	strChatIDs := make([]string, len(chatIDs))
	for i, chatID := range chatIDs {
		strChatIDs[i] = idToString(chatID)
	}

	if err := mh.pubsub.Subscribe(ctx, strChatIDs...); err != nil {
		return fmt.Errorf("repository: redis: Subscribe: %w", err)
	}

	return nil
}

func (mh *ChatPublisher) Unsubscribe(ctx context.Context, chatIDs ...int64) error {
	if chatIDs == nil {
		return nil
	}

	strChatIDs := make([]string, len(chatIDs))
	for i, chatID := range chatIDs {
		strChatIDs[i] = idToString(chatID)
	}

	if err := mh.pubsub.Unsubscribe(ctx, strChatIDs...); err != nil {
		return fmt.Errorf("repository: redis: Unsubscribe: %w", err)
	}

	return nil
}

func idToString(id int64) string {
	return strconv.FormatInt(id, 10)
}
