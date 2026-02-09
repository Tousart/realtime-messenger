package repository

import "context"

type QueuesRepository interface {
	Queues(ctx context.Context, chatID int) ([]string, error)
	AddQueueToChat(ctx context.Context, chatID int) error
	RemoveQueueFromChat(ctx context.Context, chatID int) error
}
