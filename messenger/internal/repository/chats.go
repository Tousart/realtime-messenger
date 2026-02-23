package repository

import "context"

type ChatsRepository interface {
	CreateChat(ctx context.Context)
}
