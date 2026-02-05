package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tousart/messenger/internal/models"
	"github.com/tousart/messenger/internal/repository"
)

type MessagesPublisherService struct {
	publisherRepo repository.MessagesPublisherRepository
}

func NewMessagesPublisherService(repo repository.MessagesPublisherRepository) *MessagesPublisherService {
	return &MessagesPublisherService{
		publisherRepo: repo,
	}
}

func (p *MessagesPublisherService) PublishMessage(ctx context.Context, message models.Message) error {
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
