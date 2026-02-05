package service

import (
	"context"
	"fmt"
	"router/internal/repository"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessagesConsumerService struct {
	consumerRepo repository.MessagesConsumerRepository
}

func NewMessagesConsumerService(repo repository.MessagesConsumerRepository) *MessagesConsumerService {
	return &MessagesConsumerService{
		consumerRepo: repo,
	}
}

func (c *MessagesConsumerService) ConsumeMessages(ctx context.Context) (<-chan amqp.Delivery, error) {
	messagesChannel, err := c.consumerRepo.ConsumeMessages(ctx)
	if err != nil {
		return nil, fmt.Errorf("service: ConsumeMessages: %s", err.Error())
	}

	return messagesChannel, nil
}
