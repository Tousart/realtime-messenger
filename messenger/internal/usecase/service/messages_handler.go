package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/tousart/messenger/internal/dto"
	"github.com/tousart/messenger/internal/repository"
	"github.com/tousart/messenger/internal/usecase"
)

type MessagesHandlerService struct {
	wsManager       usecase.WebSocketManager
	msgsHandlerRepo repository.MessagesHandlerRepository
}

func NewMessagesHandlerService(wsManager usecase.WebSocketManager, msgsHandlerRepo repository.MessagesHandlerRepository) *MessagesHandlerService {
	return &MessagesHandlerService{
		wsManager:       wsManager,
		msgsHandlerRepo: msgsHandlerRepo,
	}
}

func (s *MessagesHandlerService) PublishMessageToChat(ctx context.Context, input dto.SendMessageWSRequest) error {
	messagePayload, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("service: PublishMessageToChat error: %w", err)
	}

	log.Printf("сообщение опубликовано в usecase\n")

	err = s.msgsHandlerRepo.PublishMessageToChat(ctx, strconv.Itoa(input.ChatID), messagePayload)
	if err != nil {
		return fmt.Errorf("service: PublishMessageToChat error: %w", err)
	}
	return nil
}

func (s MessagesHandlerService) SubscribeToChats(ctx context.Context, chatIDs ...string) error {
	if err := s.msgsHandlerRepo.SubscribeToChats(ctx, chatIDs...); err != nil {
		return fmt.Errorf("service: SubscribeToChats error: %w", err)
	}
	return nil
}

func (s MessagesHandlerService) UnsubscribeFromChats(ctx context.Context, chatIDs ...string) error {
	if err := s.msgsHandlerRepo.UnsubscribeFromChats(ctx, chatIDs...); err != nil {
		return fmt.Errorf("service: UnsubscribeFromChats error: %w", err)
	}
	return nil
}

func (s *MessagesHandlerService) SendMessageToUsersConnections(ctx context.Context, input dto.ConsumingMessage) error {
	log.Printf("сообщение получено в usecase\n")

	if err := s.wsManager.SendMessageToUsersConnections(ctx, input); err != nil {
		return fmt.Errorf("service: SendMessageToUsersConnections error: %w", err)
	}

	// TODO: отправка сообщения в базу данных
	// message, err := domain.NewMessage(
	// 	domain.WithMessageUserID(input.UserID),
	// 	domain.WithMessageChatID(input.ChatID),
	// 	domain.WithMessageText(input.Text),
	// )

	return nil
}
