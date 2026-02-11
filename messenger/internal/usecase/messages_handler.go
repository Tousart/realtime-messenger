package usecase

import (
	"context"

	"github.com/tousart/messenger/internal/domain"
)

type MessagesHandlerService interface {
	PublishMessageToQueues(ctx context.Context, message domain.Message) error
	MessagesQueue() (domain.MessagesQueue, error)
	AddQueueToChat(ctx context.Context, chatID int) error
	RemoveQueueFromChat(ctx context.Context, chatID int) error
}
