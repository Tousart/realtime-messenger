package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tousart/messenger/internal/domain"
)

type MessagesRepository struct {
	db *sql.DB
}

func NewMessagesRepository(db *sql.DB) *MessagesRepository {
	return &MessagesRepository{
		db: db,
	}
}

func (r *MessagesRepository) Messages(ctx context.Context, chatID int64) ([]domain.Message, error) {
	const op = "repository: postgres: Messages:"

	rows, err := r.db.QueryContext(ctx,
		`SELECT message_id, user_id, chat_id, message_body, created_at
		FROM messages
		WHERE chat_id = $1
		ORDER BY created_at ASC`,
		chatID,
	)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var msg domain.Message
		if err = rows.Scan(&msg.ID, &msg.SenderID, &msg.ChatID, &msg.Text, &msg.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s %w", op, err)
		}
		messages = append(messages, msg)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	return messages, nil
}

func (r *MessagesRepository) CreateChat(ctx context.Context, chat *domain.Chat) ([]domain.ChatParticipant, error) {
	const op = "repository: postgres: CreateChat:"

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO chats (chat_id, chat_name, created_at) VALUES ($1, $2, $3)`,
		chat.ID, chat.Name, chat.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	for _, p := range chat.ChatParticipants {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO chat_user (chat_id, user_id) VALUES ($1, $2)`,
			chat.ID, p.ID,
		)
		if err != nil {
			return nil, fmt.Errorf("%s %w", op, err)
		}
	}

	rows, err := tx.QueryContext(ctx,
		`SELECT u.user_id, u.user_name
		FROM chat_user cu
		JOIN users u ON cu.user_id = u.user_id
		WHERE cu.chat_id = $1`,
		chat.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}
	defer rows.Close()

	var participants []domain.ChatParticipant
	for rows.Next() {
		var p domain.ChatParticipant
		if err = rows.Scan(&p.ID, &p.Name); err != nil {
			return nil, fmt.Errorf("%s %w", op, err)
		}
		participants = append(participants, p)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	return participants, nil
}

func (r *MessagesRepository) UsersChats(ctx context.Context, userID int64) ([]domain.ChatInfo, error) {
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

func (r *MessagesRepository) Save(ctx context.Context, msg *domain.Message) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO messages (message_id, user_id, chat_id, message_body, created_at) VALUES ($1, $2, $3, $4, $5)`,
		msg.ID, msg.SenderID, msg.ChatID, msg.Text, msg.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("repository: postgres: Save: %w", err)
	}

	return nil
}
