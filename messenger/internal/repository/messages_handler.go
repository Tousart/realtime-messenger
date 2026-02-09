package repository

import (
	"context"

	"github.com/tousart/messenger/internal/models"
)

type MessagesHandlerRepository interface {
	PublishMessageToQueues(ctx context.Context, queues []string, message models.Message) error
	MessagesQueue() (models.MessagesQueue, error)
}
