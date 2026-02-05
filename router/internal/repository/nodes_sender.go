package repository

import (
	"context"
	"router/internal/models"
)

type NodesSenderRepository interface {
	GetChatNodes(ctx context.Context, chatID string) ([]string, error)
	SendMessageToNode(ctx context.Context, nodeAddr string, message *models.Message) error
}
