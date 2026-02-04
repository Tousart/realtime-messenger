package service

import (
	"context"
	"fmt"
	"router/internal/repository"
	"strconv"
)

type RedisChatNodesSetService struct {
	setRepo repository.ChatNodesSetRepository
}

func NewRedisChatNodesSetService(repo repository.ChatNodesSetRepository) *RedisChatNodesSetService {
	return &RedisChatNodesSetService{
		setRepo: repo,
	}
}

func (r *RedisChatNodesSetService) GetChatNodes(ctx context.Context, chatID int) ([]string, error) {
	nodes, err := r.setRepo.GetChatNodes(ctx, strconv.Itoa(chatID))
	if err != nil {
		return nil, fmt.Errorf("service: GetChatNodes error: %s", err.Error())
	}
	return nodes, nil
}
