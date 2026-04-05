package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tousart/messenger/internal/domain"
	"github.com/tousart/messenger/internal/dto"
)

type MessagesUsecase struct {
	msgsRepo MessagesRepository
	chatPub  ChatPublisher
	idGen    IDGenerator
}

func NewMessagesUsecase(msgsRepo MessagesRepository, chatPub ChatPublisher, idGen IDGenerator) *MessagesUsecase {
	return &MessagesUsecase{
		msgsRepo: msgsRepo,
		chatPub:  chatPub,
		idGen:    idGen,
	}
}

func (uc *MessagesUsecase) SendMessage(ctx context.Context, input *dto.SendMessageRequest) (*dto.Message, error) {
	const op = "usecase: SendMessage:"

	text, err := domain.IsValidMessageText(input.Text)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	createdAt := timeNowUTC()

	msg := &domain.Message{
		ID:        uc.idGen.GenerateID(),
		SenderID:  input.SenderID,
		ChatID:    input.ChatID,
		Text:      text,
		CreatedAt: &createdAt,
	}

	if err = uc.msgsRepo.Save(ctx, msg); err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	// костыль
	msgBytes, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}
	// костыль

	if err = uc.chatPub.PublishMessage(ctx, input.ChatID, msgBytes); err != nil {
		return nil, fmt.Errorf("u%s %w", op, err)
	}

	return &dto.Message{
		SenderID:  msg.SenderID,
		ChatID:    msg.ChatID,
		Text:      msg.Text,
		CreatedAt: msg.CreatedAt,
	}, nil
}

func (uc *MessagesUsecase) SubscribeToChats(ctx context.Context, chatIDs ...int64) error {
	if err := uc.chatPub.Subscribe(ctx, chatIDs...); err != nil {
		return fmt.Errorf("usecase: SubscribeToChats: %w", err)
	}
	return nil
}

func (u *MessagesUsecase) UnsubscribeFromChats(ctx context.Context, chatIDs ...int64) error {
	if err := u.chatPub.Unsubscribe(ctx, chatIDs...); err != nil {
		return fmt.Errorf("usecase: UnsubscribeFromChats: %w", err)
	}
	return nil
}

func (uc *MessagesUsecase) CreateChat(ctx context.Context, input *dto.CreateChatRequest) (*dto.Chat, error) {
	const op = "usecase: CreateChat:"

	err := domain.ValidateChatName(input.ChatName)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	if len(input.ChatParticipants) == 0 {
		return nil, fmt.Errorf("%s %w", op, domain.ErrEmptyChat)
	}

	participantsIDs := make([]domain.ChatParticipant, len(input.ChatParticipants))
	for i, participant := range input.ChatParticipants {
		participantsIDs[i] = domain.ChatParticipant{
			ID: participant.ID,
		}
	}

	createdAt := timeNowUTC()

	chat := &domain.Chat{
		ID:               uc.idGen.GenerateID(),
		Name:             input.ChatName,
		ChatParticipants: participantsIDs,
		CreatedAt:        &createdAt,
	}

	participantsDB, err := uc.msgsRepo.CreateChat(ctx, chat)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	participants := make([]dto.ChatParticipant, len(participantsDB))
	for i, participant := range participantsDB {
		participants[i] = dto.ChatParticipant{
			ID:   participant.ID,
			Name: &participant.Name,
			Role: &participant.Role,
		}
	}

	return &dto.Chat{
		ID:               chat.ID,
		Name:             chat.Name,
		ChatParticipants: participants,
		CreatedAt:        chat.CreatedAt,
	}, nil
}

func (uc *MessagesUsecase) UsersChats(ctx context.Context, userID int64) ([]dto.ChatInfo, error) {
	chatsDB, err := uc.msgsRepo.UsersChats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("usecase: UsersChats: %v", err)
	}

	chats := make([]dto.ChatInfo, len(chatsDB))
	for i, chat := range chatsDB {
		chats[i] = dto.ChatInfo{
			ID:   chat.ID,
			Name: chat.Name,
		}
	}

	return chats, nil
}
