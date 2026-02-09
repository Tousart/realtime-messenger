package usecase

import (
	"context"

	"github.com/tousart/messenger/internal/models"
)

type MessagesHandlerService interface {
	PublishMessageToQueues(ctx context.Context, message models.Message) error
	MessagesQueue() (models.MessagesQueue, error)
	AddQueueToChat(ctx context.Context, chatID int) error
	RemoveQueueFromChat(ctx context.Context, chatID int) error
}
