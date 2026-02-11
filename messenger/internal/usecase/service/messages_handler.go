package service

import (
	"context"
	"fmt"

	"github.com/tousart/messenger/internal/domain"
	"github.com/tousart/messenger/internal/repository"
)

type MessagesHandlerService struct {
	msgsHandlerRepo repository.MessagesHandlerRepository
	queuesRepo      repository.QueuesRepository
}

func NewMessagesHandlerService(msgsHandlerRepo repository.MessagesHandlerRepository, queuesRepo repository.QueuesRepository) *MessagesHandlerService {
	return &MessagesHandlerService{
		msgsHandlerRepo: msgsHandlerRepo,
		queuesRepo:      queuesRepo,
	}
}

func (s *MessagesHandlerService) PublishMessageToQueues(ctx context.Context, message domain.Message) error {
	queues, err := s.queuesRepo.Queues(ctx, message.ChatID)
	if err != nil {
		return fmt.Errorf("service: PublishMessageToQueues error: %s", err.Error())
	}
	if err := s.msgsHandlerRepo.PublishMessageToQueues(ctx, queues, message); err != nil {
		return fmt.Errorf("service: PublishMessageToQueues error: %s", err.Error())
	}
	return nil
}

func (s *MessagesHandlerService) MessagesQueue() (domain.MessagesQueue, error) {
	messagesChannel, err := s.msgsHandlerRepo.MessagesQueue()
	if err != nil {
		return nil, fmt.Errorf("service: MessagesQueue: %s", err.Error())
	}
	return messagesChannel, nil
}

func (s *MessagesHandlerService) AddQueueToChat(ctx context.Context, chatID int) error {
	err := s.queuesRepo.AddQueueToChat(ctx, chatID)
	if err != nil {
		return fmt.Errorf("service: AddQueueToChat error: %s", err.Error())
	}
	return nil
}

func (s *MessagesHandlerService) RemoveQueueFromChat(ctx context.Context, chatID int) error {
	err := s.queuesRepo.RemoveQueueFromChat(ctx, chatID)
	if err != nil {
		return fmt.Errorf("service: RemoveQueueFromChat error: %s", err.Error())
	}
	return nil
}
