package repository

import "context"

type ChatNodesSetRepository interface {
	AddNodeToChat(ctx context.Context, chatID string) error
	RemoveNodeFromChat(ctx context.Context, chatID string) error
}
