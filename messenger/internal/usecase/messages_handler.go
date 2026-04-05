package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/tousart/messenger/internal/domain"
	"github.com/tousart/messenger/internal/dto"
)

type MessagesUsecase struct {
	messagesRepo MessagesHandlerRepository
	chatsRepo    ChatsRepository
}

func NewMessagesUsecase(messagesRepo MessagesHandlerRepository, chatsRepo ChatsRepository) *MessagesUsecase {
	return &MessagesUsecase{
		messagesRepo: messagesRepo,
		chatsRepo:    chatsRepo,
	}
}

func (u *MessagesUsecase) PublishMessageToChat(ctx context.Context, input *dto.SendMessageRequest) error {
	messagePayload, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("service: PublishMessageToChat error: %w", err)
	}

	err = u.messagesRepo.PublishMessageToChat(ctx, strconv.Itoa(input.ChatID), messagePayload)
	if err != nil {
		return fmt.Errorf("service: PublishMessageToChat error: %w", err)
	}
	return nil
}

func (u *MessagesUsecase) SubscribeToChats(ctx context.Context, chatIDs ...string) error {
	if err := u.messagesRepo.SubscribeToChats(ctx, chatIDs...); err != nil {
		return fmt.Errorf("service: SubscribeToChats error: %w", err)
	}
	return nil
}

func (u *MessagesUsecase) UnsubscribeFromChats(ctx context.Context, chatIDs ...string) error {
	if err := u.messagesRepo.UnsubscribeFromChats(ctx, chatIDs...); err != nil {
		return fmt.Errorf("service: UnsubscribeFromChats error: %w", err)
	}
	return nil
}

func (u *MessagesUsecase) CreateChat(ctx context.Context, input *dto.CreateChatRequest) (*dto.CreateChatResponse, error) {
	err := domain.ValidateChatName(input.ChatName)
	if err != nil {
		return nil, fmt.Errorf("service: CreateChat error: %w", domain.ErrInvalidRequest)
	}

	if len(input.ChatParticipants) == 0 {
		return nil, fmt.Errorf("service: CreateChat error: %w", domain.ErrEmptyChat)
	}

	chat := &domain.Chat{
		Name: input.ChatName,
	}

	userNames := make([]string, len(input.ChatParticipants))
	for i, participant := range input.ChatParticipants {
		userNames[i] = participant.UserName
	}

	createdChat, err := u.chatsRepo.CreateChat(ctx, chat, userNames...)
	if err != nil {
		return nil, fmt.Errorf("service: CreateChat error: %w", err)
	}

	participants := make([]dto.ChatParticipantResponse, len(createdChat.ChatParticipants))
	for i, participant := range createdChat.ChatParticipants {
		participants[i] = dto.ChatParticipantResponse{
			UserID:   participant.UserID,
			UserName: participant.UserName,
			Role:     participant.Role,
		}
	}
	return &dto.CreateChatResponse{
		ChatID:           createdChat.ID,
		ChatName:         createdChat.Name,
		ChatParticipants: participants,
	}, nil
}

func (u *MessagesUsecase) UsersChats(ctx context.Context, userID int) ([]dto.Chat, error) {
	chatsDB, err := u.chatsRepo.UsersChats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("usecase: UsersChats: %v", err)
	}

	chats := make([]dto.Chat, len(chatsDB))
	for i, chat := range chatsDB {
		chats[i] = dto.Chat{
			ID:   chat.ID,
			Name: chat.Name,
		}
	}

	return chats, nil
}
