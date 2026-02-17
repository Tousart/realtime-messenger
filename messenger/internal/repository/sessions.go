package repository

import (
	"context"

	"github.com/tousart/messenger/internal/domain"
)

type SessionsRepository interface {
	SessionData(ctx context.Context, sessionID string) (*domain.User, error)
	GenerateSessionID(ctx context.Context, user *domain.User) (string, error)
}
