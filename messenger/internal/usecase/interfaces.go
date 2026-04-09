package usecase

import (
	"context"

	"github.com/tousart/messenger/internal/domain"
)

type MessagesRepository interface {
	Messages(ctx context.Context, chatID int64) ([]domain.Message, error)
	CreateChat(ctx context.Context, chat *domain.Chat) ([]domain.ChatParticipant, error)
	UsersChats(ctx context.Context, userID int64) ([]domain.ChatInfo, error)
	Save(ctx context.Context, msg *domain.Message) error
}

type UsersRepository interface {
	Create(ctx context.Context, user *domain.User) error
	User(ctx context.Context, name string) (*domain.User, error)
}

type ChatPublisher interface {
	PublishMessage(ctx context.Context, msg *domain.Message) error
	Subscribe(ctx context.Context, chatIDs ...int64) error
	Unsubscribe(ctx context.Context, chatIDs ...int64) error
}

type SessionsRepository interface {
	GenerateSessionID(ctx context.Context, payload *domain.SessionPayload) (string, error)
	Payload(ctx context.Context, sessionID string) (*domain.SessionPayload, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) bool
}

type IDGenerator interface {
	GenerateID() int64
}
