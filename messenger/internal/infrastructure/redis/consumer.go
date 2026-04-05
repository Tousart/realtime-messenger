package redis

import (
	"context"
	"encoding/json"
	"log"

	rdb "github.com/redis/go-redis/v9"
	"github.com/tousart/messenger/internal/dto"
)

type WebsocketManager interface {
	SendMessageToUsersConnections(ctx context.Context, msg dto.ConsumingMessage) error
}

type Consumer struct {
	wsManager WebsocketManager
	pubsub    *rdb.PubSub
}

func NewRedisConsumer(wsManager WebsocketManager, pubsub *rdb.PubSub) *Consumer {
	return &Consumer{
		wsManager: wsManager,
		pubsub:    pubsub,
	}
}

func (c *Consumer) ConsumeMessages(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-c.pubsub.Channel():
			var consumingMessage dto.ConsumingMessage
			if err := json.Unmarshal([]byte(msg.Payload), &consumingMessage); err != nil {
				log.Printf("infrastructure: ConsumeMessages: failed to unmarshal message: %v\n", err)
				continue
			}

			if err := c.wsManager.SendMessageToUsersConnections(ctx, consumingMessage); err != nil {
				log.Printf("infrastructure: ConsumeMessages: failed to send message to users connections: %v\n", err)
				continue
			}
		}
	}
}
