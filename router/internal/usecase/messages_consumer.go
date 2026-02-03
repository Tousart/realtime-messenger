package usecase

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessagesConsumerService interface {
	ConsumeMessages(ctx context.Context) (<-chan amqp.Delivery, error)
}
