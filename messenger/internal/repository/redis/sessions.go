package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tousart/messenger/internal/domain"

	rdb "github.com/redis/go-redis/v9"
)

const (
	SessionExpirationHours = 6
)

type SessionsRepository struct {
	client *rdb.Client
}

func NewSessionsRepository(client *rdb.Client) *SessionsRepository {
	return &SessionsRepository{
		client: client,
	}
}

func (r *SessionsRepository) Payload(ctx context.Context, sessionID string) ([]byte, error) {
	const op = "repository: redis: Payload:"

	payload, err := r.client.Get(ctx, sessionID).Bytes()
	if err != nil {
		if payload == nil {
			return nil, fmt.Errorf("%s %w", op, domain.ErrSessionIDNotExists)
		}
		return nil, fmt.Errorf("%s %w", op, err)
	}

	return payload, nil
}

func (r *SessionsRepository) GenerateSessionID(ctx context.Context, payload []byte) (string, error) {
	sessionID := uuid.New().String()
	if err := r.client.Set(ctx, sessionID, payload, SessionExpirationHours*time.Hour).Err(); err != nil {
		return "", fmt.Errorf("repository: redis: GenerateSessionID: %w", err)
	}
	return sessionID, nil
}
