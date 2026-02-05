package repository

import "context"

type NodesReceiverRepository interface {
	AddNodeToChat(ctx context.Context, chatID string) error
	RemoveNodeFromChat(ctx context.Context, chatID string) error
}
