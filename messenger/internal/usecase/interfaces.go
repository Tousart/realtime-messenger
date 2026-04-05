package usecase

import (
	"context"

	"github.com/tousart/messenger/internal/domain"
)

type MessagesRepository interface {
	CreateChat(ctx context.Context, chat *domain.Chat) ([]domain.ChatParticipant, error)
	UsersChats(ctx context.Context, userID int64) ([]domain.ChatInfo, error)
	Save(ctx context.Context, msg *domain.Message) error
}

type UsersRepository interface {
	Create(ctx context.Context, user *domain.User) error
	User(ctx context.Context, name string) (*domain.User, error)
}

type ChatPublisher interface {
	PublishMessage(ctx context.Context, chatID int64, msgBytes []byte) error
	Subscribe(ctx context.Context, chatIDs ...int64) error
	Unsubscribe(ctx context.Context, chatIDs ...int64) error
}

type SessionsRepository interface {
	GenerateSessionID(ctx context.Context, payload []byte) (string, error)
	Payload(ctx context.Context, sessionID string) ([]byte, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) bool
}

type IDGenerator interface {
	GenerateID() int64
}
