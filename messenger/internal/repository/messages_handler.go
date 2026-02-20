package repository

import (
	"context"
)

type MessagesHandlerRepository interface {
	PublishMessageToChat(ctx context.Context, chatID string, messagePayload []byte) error
	SubscribeToChats(ctx context.Context, chatIDs ...string) error
	UnsubscribeFromChats(ctx context.Context, chatIDs ...string) error
}
