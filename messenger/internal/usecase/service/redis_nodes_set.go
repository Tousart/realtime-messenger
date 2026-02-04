package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/tousart/messenger/internal/repository"
)

type RedisChatNodesSetService struct {
	setRepo repository.ChatNodesSetRepository
}

func NewRedisChatNodesSetService(repo repository.ChatNodesSetRepository) *RedisChatNodesSetService {
	return &RedisChatNodesSetService{
		setRepo: repo,
	}
}

func (r *RedisChatNodesSetService) AddNodeToChat(ctx context.Context, chatID int) error {
	err := r.setRepo.AddNodeToChat(ctx, strconv.Itoa(chatID))
	if err != nil {
		return fmt.Errorf("service: AddNodeToChat error: %s", err.Error())
	}
	return nil
}

func (r *RedisChatNodesSetService) RemoveNodeFromChat(ctx context.Context, chatID int) error {
	err := r.setRepo.RemoveNodeFromChat(ctx, strconv.Itoa(chatID))
	if err != nil {
		return fmt.Errorf("service: RemoveNodeFromChat error: %s", err.Error())
	}
	return nil
}
