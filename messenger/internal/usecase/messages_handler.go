package usecase

import (
	"context"

	"github.com/tousart/messenger/internal/domain"
	"github.com/tousart/messenger/internal/dto"
)

type MessagesHandlerService interface {
	PublishMessageToQueues(ctx context.Context, message *dto.SendMessageWSRequest) error
	MessagesQueue() (domain.MessagesQueue, error)
	AddQueueToChat(ctx context.Context, input dto.ChatWSRequest) error
	RemoveQueueFromChat(ctx context.Context, input dto.ChatWSRequest) error
}
