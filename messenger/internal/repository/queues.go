package repository

import (
	"context"

	"github.com/tousart/messenger/internal/domain"
)

type QueuesRepository interface {
	Queues(ctx context.Context, chat *domain.Chat) ([]string, error)
	AddQueueToChat(ctx context.Context, chat *domain.Chat) error
	RemoveQueueFromChat(ctx context.Context, chat *domain.Chat) error
}
