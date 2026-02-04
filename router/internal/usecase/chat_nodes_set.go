package usecase

import "context"

type ChatNodesSetService interface {
	GetChatNodes(ctx context.Context, chatID int) ([]string, error)
}
