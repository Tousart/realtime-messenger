package service

import (
	"context"
	"fmt"
	"router/internal/repository"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumerService struct {
	consumerRepo repository.MessagesConsumerRepository
}

func NewRabbitMQConsumerService(repo repository.MessagesConsumerRepository) *RabbitMQConsumerService {
	return &RabbitMQConsumerService{
		consumerRepo: repo,
	}
}

func (c *RabbitMQConsumerService) ConsumeMessages(ctx context.Context) (<-chan amqp.Delivery, error) {
	messagesChannel, err := c.consumerRepo.ConsumeMessages(ctx)
	if err != nil {
		return nil, fmt.Errorf("service: ConsumeMessages: %s", err.Error())
	}

	return messagesChannel, nil
}
