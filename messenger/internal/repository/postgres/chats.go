package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/tousart/messenger/internal/domain"
)

type ChatsRepository struct {
	db *sql.DB
}

func NewChatsRepository(db *sql.DB) *ChatsRepository {
	return &ChatsRepository{
		db: db,
	}
}

func (r *ChatsRepository) CreateChat(ctx context.Context, chat *domain.Chat, userNames ...string) (*domain.Chat, error) {
	const op = "postgres: CreateChat:"

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, "SELECT user_id, user_name FROM users WHERE user_name = ANY($1)", userNames)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}
	defer rows.Close()

	participants := make([]domain.ChatParticipant, len(userNames))
	idx := 0
	for rows.Next() {
		var participant domain.ChatParticipant
		if err = rows.Scan(&participant.UserID, participant.UserName); err != nil {
			return nil, fmt.Errorf("%s %w", op, err)
		}
		participants[idx] = participant
		idx++
	}
	if idx != len(userNames) {
		return nil, fmt.Errorf("%s %w", op, domain.ErrUserNotFound)
	}

	var chatID int
	if err = r.db.QueryRowContext(ctx, "INSERT INTO chats (chat_name, created_at, updated_at) VALUES ($1, NOW(), NOW())", chat.Name).Scan(&chatID); err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	reflectionsUserChat := make([]string, len(participants))
	for i, p := range participants {
		reflectionsUserChat[i] = fmt.Sprintf("(%d,%d)", chatID, p.UserID)
	}
	reflectionArgs := strings.Join(reflectionsUserChat, ",")
	if _, err = r.db.ExecContext(ctx, "INSERT INTO chat_user (chat_id, user_id) VALUES $1", reflectionArgs); err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}
	return &domain.Chat{
		ID:               chatID,
		Name:             chat.Name,
		ChatParticipants: participants,
	}, nil
}

func (r *ChatsRepository) UsersChats(ctx context.Context, userID int) ([]domain.ChatInfo, error) {
	const op = "postgres: UsersChats:"

	rows, err := r.db.QueryContext(ctx,
		`SELECT c.chat_id, c.chat_name
		FROM chats c
		JOIN chat_user cu ON c.chat_id = cu.chat_id
		WHERE cu.user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}
	defer rows.Close()

	var chats []domain.ChatInfo
	for rows.Next() {
		var chat domain.ChatInfo
		if err = rows.Scan(&chat.ID, &chat.Name); err != nil {
			return nil, fmt.Errorf("%s %w", op, err)
		}
		chats = append(chats, chat)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	return chats, nil
}
