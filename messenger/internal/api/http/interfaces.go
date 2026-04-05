package api

import (
	"context"
	"net/http"

	"github.com/tousart/messenger/internal/dto"
)

type MessagesUsecase interface {
	UsersChats(ctx context.Context, userID int64) ([]dto.ChatInfo, error)
}

type UsersUsecase interface {
	Register(ctx context.Context, input *dto.RegisterRequest) (*dto.User, error)
	Login(ctx context.Context, input *dto.LoginRequest) (*dto.SessionID, error)
	ValidateSessionID(ctx context.Context, sessionID string) ([]byte, error)
}

type WebSocketUpgrader interface {
	UpgradeConnectionForUser(w http.ResponseWriter, r *http.Request, responseHeader http.Header, payload *dto.UserPayload) error
}
