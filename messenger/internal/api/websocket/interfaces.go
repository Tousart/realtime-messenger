package websocket

import (
	"context"

	"github.com/tousart/messenger/internal/dto"
)

type MessagesUsecase interface {
	SubscribeToChats(ctx context.Context, chatIDs ...string) error
	UnsubscribeFromChats(ctx context.Context, chatIDs ...string) error
	PublishMessageToChat(ctx context.Context, input *dto.SendMessageRequest) error
	CreateChat(ctx context.Context, input *dto.CreateChatRequest) (*dto.CreateChatResponse, error)
}
