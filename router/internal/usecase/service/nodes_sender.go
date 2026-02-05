package service

import (
	"context"
	"fmt"
	"router/internal/models"
	"router/internal/repository"
	"strconv"
	"strings"
)

type NodesSenderService struct {
	senderRepo repository.NodesSenderRepository
}

func NewNodesSenderService(repo repository.NodesSenderRepository) *NodesSenderService {
	return &NodesSenderService{
		senderRepo: repo,
	}
}

func (ns *NodesSenderService) GetChatNodes(ctx context.Context, chatID int) ([]string, error) {
	nodes, err := ns.senderRepo.GetChatNodes(ctx, strconv.Itoa(chatID))
	if err != nil {
		return nil, fmt.Errorf("service: GetChatNodes error: %s", err.Error())
	}
	return nodes, nil
}

func (ns *NodesSenderService) SendMessageToNode(ctx context.Context, nodeAddr string, message *models.Message) error {
	if !strings.HasPrefix(nodeAddr, "http://") {
		nodeAddr = "http://" + nodeAddr
	}
	if !strings.HasSuffix(nodeAddr, "/messages") {
		nodeAddr = nodeAddr + "/messages"
	}
	err := ns.senderRepo.SendMessageToNode(ctx, nodeAddr, message)
	if err != nil {
		return fmt.Errorf("service: SendMessageToNode error: %s", err.Error())
	}
	return nil
}
