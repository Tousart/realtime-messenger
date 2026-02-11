package usecase

import (
	"context"

	"github.com/tousart/messenger/internal/domain"
)

type UsersService interface {
	RegisterUser(ctx context.Context, data *domain.RegisterRequest) error
	LoginUser(ctx context.Context, data *domain.LoginRequest) error
}
