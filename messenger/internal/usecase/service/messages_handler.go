package service

import (
	"context"
	"fmt"

	"github.com/tousart/messenger/internal/domain"
	"github.com/tousart/messenger/internal/dto"
	"github.com/tousart/messenger/internal/repository"
	"github.com/tousart/messenger/internal/usecase"
)

type MessagesHandlerService struct {
	wsManager       usecase.WebSocketManager
	msgsHandlerRepo repository.MessagesHandlerRepository
	queuesRepo      repository.QueuesRepository
}

func NewMessagesHandlerService(wsManager usecase.WebSocketManager, msgsHandlerRepo repository.MessagesHandlerRepository, queuesRepo repository.QueuesRepository) *MessagesHandlerService {
	return &MessagesHandlerService{
		wsManager:       wsManager,
		msgsHandlerRepo: msgsHandlerRepo,
		queuesRepo:      queuesRepo,
	}
}

func (s *MessagesHandlerService) PublishMessageToQueues(ctx context.Context, input dto.SendMessageWSRequest) error {
	message, err := domain.NewMessage(
		domain.WithMessageUserID(input.UserID),
		domain.WithMessageChatID(input.ChatID),
		domain.WithMessageText(input.Text),
	)
	if err != nil {
		return fmt.Errorf("service: PublishMessageToQueues error: %w", err)
	}

	chat, err := domain.NewChat(domain.WithChatChatID(input.ChatID))
	if err != nil {
		return fmt.Errorf("service: PublishMessageToQueues error: %w", err)
	}

	queues, err := s.queuesRepo.Queues(ctx, chat)
	if err != nil {
		return fmt.Errorf("service: PublishMessageToQueues error: %w", err)
	}

	if err := s.msgsHandlerRepo.PublishMessageToQueues(ctx, queues, message); err != nil {
		return fmt.Errorf("service: PublishMessageToQueues error: %w", err)
	}
	return nil
}

func (s *MessagesHandlerService) SendMessageToUsersConnections(ctx context.Context, input dto.ConsumingMessage) error {
	if err := s.wsManager.SendMessageToUsersConnections(ctx, input); err != nil {
		return fmt.Errorf("service: SendMessageToUsersConnections error: %w", err)
	}

	// TODO: отправка сообщения в базу данных

	return nil
}

func (s *MessagesHandlerService) AddQueueToChat(ctx context.Context, input dto.ChatWSRequest) error {
	chat, err := domain.NewChat(domain.WithChatChatID(input.ChatID))
	if err != nil {
		return fmt.Errorf("service: PublishMessageToQueues error: %w", err)
	}

	if err = s.queuesRepo.AddQueueToChat(ctx, chat); err != nil {
		return fmt.Errorf("service: AddQueueToChat error: %w", err)
	}
	return nil
}

func (s *MessagesHandlerService) RemoveQueueFromChat(ctx context.Context, input dto.ChatWSRequest) error {
	chat, err := domain.NewChat(domain.WithChatChatID(input.ChatID))
	if err != nil {
		return fmt.Errorf("service: PublishMessageToQueues error: %w", err)
	}

	if err = s.queuesRepo.RemoveQueueFromChat(ctx, chat); err != nil {
		return fmt.Errorf("service: RemoveQueueFromChat error: %w", err)
	}
	return nil
}
