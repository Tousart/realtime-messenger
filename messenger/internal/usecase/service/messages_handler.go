package service

import (
	"context"
	"fmt"

	"github.com/tousart/messenger/internal/domain"
	"github.com/tousart/messenger/internal/dto"
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

func (s *MessagesHandlerService) PublishMessageToQueues(ctx context.Context, input *dto.SendMessageWSRequest) error {
	message, err := domain.NewMessage(
		domain.WithMessageUserID(input.UserID),
		domain.WithMessageChatID(input.ChatID),
		domain.WithMessageText(input.Text),
	)
	if err != nil {
		return fmt.Errorf("service: PublishMessageToQueues error: %s", err.Error())
	}

	chat, err := domain.NewChat(domain.WithChatChatID(input.ChatID))
	if err != nil {
		return fmt.Errorf("service: PublishMessageToQueues error: %s", err.Error())
	}

	queues, err := s.queuesRepo.Queues(ctx, chat)
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

func (s *MessagesHandlerService) AddQueueToChat(ctx context.Context, input dto.ChatWSRequest) error {
	chat, err := domain.NewChat(domain.WithChatChatID(input.ChatID))
	if err != nil {
		return fmt.Errorf("service: PublishMessageToQueues error: %s", err.Error())
	}

	if err = s.queuesRepo.AddQueueToChat(ctx, chat); err != nil {
		return fmt.Errorf("service: AddQueueToChat error: %s", err.Error())
	}
	return nil
}

func (s *MessagesHandlerService) RemoveQueueFromChat(ctx context.Context, input dto.ChatWSRequest) error {
	chat, err := domain.NewChat(domain.WithChatChatID(input.ChatID))
	if err != nil {
		return fmt.Errorf("service: PublishMessageToQueues error: %s", err.Error())
	}

	if err = s.queuesRepo.RemoveQueueFromChat(ctx, chat); err != nil {
		return fmt.Errorf("service: RemoveQueueFromChat error: %s", err.Error())
	}
	return nil
}
