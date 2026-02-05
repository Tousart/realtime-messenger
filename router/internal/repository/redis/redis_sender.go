package redis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"router/internal/models"
	"time"

	rdb "github.com/redis/go-redis/v9"
)

type RedisNodesSenderRepository struct {
	client *rdb.Client
}

func NewRedisNodesSenderRepository(redisAddr string) *RedisNodesSenderRepository {
	client := rdb.NewClient(&rdb.Options{
		Addr: redisAddr,
	})
	return &RedisNodesSenderRepository{
		client: client,
	}
}

func (r *RedisNodesSenderRepository) GetChatNodes(ctx context.Context, chatID string) ([]string, error) {
	nodes, err := r.client.SMembers(ctx, chatID).Result()
	if err != nil {
		return nil, fmt.Errorf("redis: GetChatNodes error: %s", err.Error())
	}
	return nodes, nil
}

func (r *RedisNodesSenderRepository) SendMessageToNode(ctx context.Context, nodeAddr string, message *models.Message) error {
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("repository: SendMessageToNode error: %s", err.Error())
	}

	req, err := http.NewRequestWithContext(ctx, "POST", nodeAddr, bytes.NewBuffer(msgBytes))
	if err != nil {
		return fmt.Errorf("repository: SendMessageToNode error: %s", err.Error())
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("repository: SendMessageToNode error: %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("repository: SendMessageToNode error: %d", resp.StatusCode)
	}
	return nil
}
