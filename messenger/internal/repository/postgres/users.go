package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/tousart/messenger/internal/domain"
)

type UsersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) (*UsersRepository, error) {
	return &UsersRepository{
		db: db,
	}, nil
}

func (r *UsersRepository) Create(ctx context.Context, user *domain.User) error {
	const op = "repository: postgres: Create:"

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO users (user_id, user_name, password, created_at) VALUES ($1, $2, $3, $4)`,
		user.ID, user.Name, user.Password, user.CreatedAt,
	)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%s %w", op, domain.ErrUserAlreadyExists)
		}
		return fmt.Errorf("%s %w", op, err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("%s %w", op, err)
	}
	return nil
}

func (r *UsersRepository) User(ctx context.Context, name string) (*domain.User, error) {
	const op = "repository: postgres: User:"

	var user domain.User
	err := r.db.QueryRowContext(ctx,
		`SELECT user_id, user_name, password, created_at FROM users WHERE user_name = $1`,
		name,
	).Scan(&user.ID, &user.Name, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s %w", op, domain.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s %w", op, err)
	}

	return &user, nil
}
