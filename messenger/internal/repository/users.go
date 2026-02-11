package repository

import (
	"context"

	"github.com/tousart/messenger/internal/domain"
)

type UsersRepository interface {
	RegisterUser(ctx context.Context, user *domain.User) error
	User(ctx context.Context, userName string) (*domain.User, error)
}
