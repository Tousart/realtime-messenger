package usecase

import (
	"context"

	"github.com/tousart/messenger/internal/dto"
)

type UsersService interface {
	RegisterUser(ctx context.Context, input dto.RegisterUserRequest) (*dto.RegisterUserResponse, error)
	LoginUser(ctx context.Context, input dto.LoginUserRequest) (*dto.LoginUserResponse, error)
	ValidateSessionID(ctx context.Context, sessionID string) (*dto.UserPayload, error)
}
