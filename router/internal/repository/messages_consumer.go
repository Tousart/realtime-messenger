package repository

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessagesConsumerRepository interface {
	ConsumeMessages(ctx context.Context) (<-chan amqp.Delivery, error)
}
