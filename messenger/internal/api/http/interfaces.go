package api

import (
	"context"
	"net/http"

	"github.com/tousart/messenger/internal/dto"
)

type MessagesUsecase interface {
	UsersChats(ctx context.Context, userID int) ([]dto.Chat, error)
}

type UsersUsecase interface {
	RegisterUser(ctx context.Context, input dto.RegisterUserRequest) (*dto.RegisterUserResponse, error)
	LoginUser(ctx context.Context, input dto.LoginUserRequest) (*dto.LoginUserResponse, error)
	ValidateSessionID(ctx context.Context, sessionID string) (*dto.SessionPayload, error)
}

type WebSocketUpgrader interface {
	UpgradeConnectionForUser(w http.ResponseWriter, r *http.Request, responseHeader http.Header, payload *dto.UserPayload) error
}
