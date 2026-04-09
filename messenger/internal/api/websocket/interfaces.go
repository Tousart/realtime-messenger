package websocket

import (
	"context"

	"github.com/tousart/messenger/internal/dto"
)

type MessagesUsecase interface {
	// messages
	Messages(ctx context.Context, chatID int64) (*dto.Messages, error)
	SendMessage(ctx context.Context, input *dto.SendMessageRequest) (*dto.Message, error)

	// chats
	SubscribeToChats(ctx context.Context, chatIDs ...int64) error
	UnsubscribeFromChats(ctx context.Context, chatIDs ...int64) error
	CreateChat(ctx context.Context, input *dto.CreateChatRequest) (*dto.Chat, error)
}
