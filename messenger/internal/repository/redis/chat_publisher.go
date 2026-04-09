package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tousart/messenger/internal/domain"
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

// соответствует dto.Message
type messagePayload struct {
	ID        int64      `json:"message_id,string"`
	SenderID  int64      `json:"sender_id,string"`
	ChatID    int64      `json:"chat_id,string"`
	Text      string     `json:"text"`
	CreatedAt *time.Time `json:"created_at"`
}

func (p *ChatPublisher) PublishMessage(ctx context.Context, msg *domain.Message) error {
	const op = "repository: redis: PublishMessage:"

	payload := messagePayload{
		ID:        msg.ID,
		SenderID:  msg.SenderID,
		ChatID:    msg.ChatID,
		Text:      msg.Text,
		CreatedAt: msg.CreatedAt,
	}

	msgBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}

	if err := p.client.Publish(ctx, idToString(msg.ChatID), msgBytes).Err(); err != nil {
		return fmt.Errorf("%s %w", op, err)
	}
	return nil
}

func (p *ChatPublisher) Subscribe(ctx context.Context, chatIDs ...int64) error {
	if chatIDs == nil {
		return nil
	}

	strChatIDs := make([]string, len(chatIDs))
	for i, chatID := range chatIDs {
		strChatIDs[i] = idToString(chatID)
	}

	if err := p.pubsub.Subscribe(ctx, strChatIDs...); err != nil {
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
