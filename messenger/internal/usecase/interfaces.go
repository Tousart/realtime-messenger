package usecase

import (
	"context"

	"github.com/tousart/messenger/internal/domain"
)

type ChatsRepository interface {
	CreateChat(ctx context.Context, chat *domain.Chat, userNames ...string) (*domain.Chat, error)
	UsersChats(ctx context.Context, userID int) ([]domain.ChatInfo, error)
}

type MessagesHandlerRepository interface {
	PublishMessageToChat(ctx context.Context, chatID string, messagePayload []byte) error
	SubscribeToChats(ctx context.Context, chatIDs ...string) error
	UnsubscribeFromChats(ctx context.Context, chatIDs ...string) error
}

type SessionsRepository interface {
	SessionIDPayload(ctx context.Context, sessionID string) (*domain.User, error)
	GenerateSessionID(ctx context.Context, user *domain.User) (string, error)
}

type UsersRepository interface {
	RegisterUser(ctx context.Context, user *domain.User) (int, error)
	User(ctx context.Context, userName string) (*domain.User, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) bool
}
