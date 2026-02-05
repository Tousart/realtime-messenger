package usecase

import (
	"context"
	"router/internal/models"
)

type NodesSenderService interface {
	GetChatNodes(ctx context.Context, chatID int) ([]string, error)
	SendMessageToNode(ctx context.Context, nodeAddr string, message *models.Message) error
}
