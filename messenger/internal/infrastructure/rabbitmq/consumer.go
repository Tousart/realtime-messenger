package rabbitmq

import (
	"context"
	"log"

	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tousart/messenger/internal/dto"
	"github.com/tousart/messenger/internal/usecase"
)

type RabbitMQConsumer struct {
	msgsHandlerService usecase.MessagesHandlerService
	messagesQueue      <-chan amqp.Delivery
}

func NewRabbitMQConsumer(msgsHandlerService usecase.MessagesHandlerService, queue <-chan amqp.Delivery) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		msgsHandlerService: msgsHandlerService,
		messagesQueue:      queue,
	}
}

func (c *RabbitMQConsumer) ConsumeMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-c.messagesQueue:
			var consumingMessage dto.ConsumingMessage
			if err := json.Unmarshal(msg.Body, &consumingMessage); err != nil {
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
