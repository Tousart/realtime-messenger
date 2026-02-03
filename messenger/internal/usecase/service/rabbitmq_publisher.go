package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tousart/messenger/internal/models"
	"github.com/tousart/messenger/internal/repository"
)

type RabbitMQPublisherService struct {
	publisherRepo repository.MessagesPublisherRepository
}

func NewRabbitMQPublisherService(repo repository.MessagesPublisherRepository) *RabbitMQPublisherService {
	return &RabbitMQPublisherService{
		publisherRepo: repo,
	}
}

func (p *RabbitMQPublisherService) PublishMessage(ctx context.Context, message models.Message) error {
	messageBytes, err := json.Marshal(message)

	if err != nil {
		return fmt.Errorf("service: PublishMessage error: %s", err.Error())
	}

	err = p.publisherRepo.PublishMessage(ctx, messageBytes)
	if err != nil {
		return fmt.Errorf("service: PublishMessage error: %s", err.Error())
	}

	return nil
}
