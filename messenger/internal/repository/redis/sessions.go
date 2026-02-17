package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tousart/messenger/internal/domain"

	rdb "github.com/redis/go-redis/v9"
)

const (
	SESSION_EXPIRATION_HOURS = 6
)

type RedisSessionsRepository struct {
	client *rdb.Client
}

func NewRedisSessionsRepository(client *rdb.Client) *RedisSessionsRepository {
	return &RedisSessionsRepository{
		client: client,
	}
}

func (r *RedisSessionsRepository) SessionData(ctx context.Context, sessionID string) (*domain.User, error) {
	data, err := r.client.Get(ctx, sessionID).Bytes()
	if err != nil {
		return nil, fmt.Errorf("redis: SessionID: %w", err)
	}

	var user domain.User
	if err := json.Unmarshal(data, &sessionID); err != nil {
		return nil, fmt.Errorf("redis: SessionID: %w", err)
	}
	return &user, nil
}

func (r *RedisSessionsRepository) GenerateSessionID(ctx context.Context, user *domain.User) (string, error) {
	sessionID := uuid.New().String()
	data, err := json.Marshal(*user)
	if err != nil {
		return "", fmt.Errorf("redis: GenerateSessionID: %w", err)
	}

	if err := r.client.Set(ctx, sessionID, data, 6*time.Hour).Err(); err != nil {
		return "", fmt.Errorf("redis: GenerateSessionID: %w", err)
	}
	return sessionID, nil
}
