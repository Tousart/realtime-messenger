package usecase

import "context"

type NodesReceiverService interface {
	AddNodeToChat(ctx context.Context, chatID int) error
	RemoveNodeFromChat(ctx context.Context, chatID int) error
}
