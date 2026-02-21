package redis

import (
	"context"
	"encoding/json"
	"log"

	rdb "github.com/redis/go-redis/v9"
	"github.com/tousart/messenger/internal/dto"
	"github.com/tousart/messenger/internal/usecase"
)

type RedisConsumer struct {
	msgsHandlerService usecase.MessagesHandlerService
	pubsub             *rdb.PubSub
}

func NewRedisConsumer(msgsHandlerService usecase.MessagesHandlerService, pubsub *rdb.PubSub) *RedisConsumer {
	return &RedisConsumer{
		msgsHandlerService: msgsHandlerService,
		pubsub:             pubsub,
	}
}

func (c *RedisConsumer) ConsumeMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-c.pubsub.Channel():
			var consumingMessage dto.ConsumingMessage
			if err := json.Unmarshal([]byte(msg.Payload), &consumingMessage); err != nil {
				log.Printf("infrastructure: ConsumeMessages: failed to unmarshal message: %v\n", err)
				continue
			}

			if err := c.msgsHandlerService.SendMessageToUsersConnections(context.Background(), consumingMessage); err != nil {
				log.Printf("infrastructure: ConsumeMessages: failed to send message to users connections: %v\n", err)
				continue
			}
		}
	}
}
