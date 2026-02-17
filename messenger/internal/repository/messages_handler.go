package repository

import (
	"context"

	"github.com/tousart/messenger/internal/domain"
)

type MessagesHandlerRepository interface {
	PublishMessageToQueues(ctx context.Context, queues []string, message *domain.Message) error
	MessagesQueue() (domain.MessagesQueue, error)
}
