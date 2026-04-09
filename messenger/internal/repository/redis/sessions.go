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

// соответствует dto.SessionPayload
type sessionPayload struct {
	UserID   int64  `json:"user_id,string"`
	UserName string `json:"user_name"`
}

func (r *SessionsRepository) Payload(ctx context.Context, sessionID string) (*domain.SessionPayload, error) {
	const op = "repository: redis: Payload:"

	payloadBytes, err := r.client.Get(ctx, sessionID).Bytes()
	if err != nil {
		if payloadBytes == nil {
			return nil, fmt.Errorf("%s %w", op, domain.ErrSessionIDNotExists)
		}
		return nil, fmt.Errorf("%s %w", op, err)
	}

	var sp domain.SessionPayload
	if err = json.Unmarshal(payloadBytes, &sp); err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	payload := &domain.SessionPayload{
		UserID:   sp.UserID,
		UserName: sp.UserName,
	}

	return payload, nil
}

func (r *SessionsRepository) GenerateSessionID(ctx context.Context, payload *domain.SessionPayload) (string, error) {
	const op = "repository: redis: GenerateSessionID:"

	sp := sessionPayload{
		UserID:   payload.UserID,
		UserName: payload.UserName,
	}

	payloadBytes, err := json.Marshal(sp)
	if err != nil {
		return "", fmt.Errorf("%s %w", op, err)
	}

	sessionID := uuid.New().String()
	if err := r.client.Set(ctx, sessionID, payloadBytes, SessionExpirationHours*time.Hour).Err(); err != nil {
		return "", fmt.Errorf("%s %w", op, err)
	}

	return sessionID, nil
}
