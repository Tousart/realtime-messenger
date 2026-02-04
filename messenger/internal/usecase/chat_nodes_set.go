package usecase

import "context"

type ChatNodesSetService interface {
	AddNodeToChat(ctx context.Context, chatID int) error
	RemoveNodeFromChat(ctx context.Context, chatID int) error
}
