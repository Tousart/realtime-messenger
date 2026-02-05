package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/tousart/messenger/internal/repository"
)

type NodesReceiverService struct {
	receiverRepo repository.NodesReceiverRepository
}

func NewNodesReceiverService(repo repository.NodesReceiverRepository) *NodesReceiverService {
	return &NodesReceiverService{
		receiverRepo: repo,
	}
}

func (r *NodesReceiverService) AddNodeToChat(ctx context.Context, chatID int) error {
	err := r.receiverRepo.AddNodeToChat(ctx, strconv.Itoa(chatID))
	if err != nil {
		return fmt.Errorf("service: AddNodeToChat error: %s", err.Error())
	}
	return nil
}

func (r *NodesReceiverService) RemoveNodeFromChat(ctx context.Context, chatID int) error {
	err := r.receiverRepo.RemoveNodeFromChat(ctx, strconv.Itoa(chatID))
	if err != nil {
		return fmt.Errorf("service: RemoveNodeFromChat error: %s", err.Error())
	}
	return nil
}
