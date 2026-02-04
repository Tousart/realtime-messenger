package repository

import "context"

type ChatNodesSetRepository interface {
	GetChatNodes(ctx context.Context, chatID string) ([]string, error)
}
