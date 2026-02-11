package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/tousart/messenger/internal/domain"
	"github.com/tousart/messenger/pkg"
)

type PSQLUsersRepository struct {
	db *sql.DB
}

func NewPSQLUsersRepository(psqlAddress string) (*PSQLUsersRepository, error) {
	db, err := pkg.ConnectToPSQL(psqlAddress)
	if err != nil {
		return nil, fmt.Errorf("postgres: NewPSQLUsersRepository: %w", err)
	}
	return &PSQLUsersRepository{
		db: db,
	}, nil
}

func (r *PSQLUsersRepository) RegisterUser(ctx context.Context, user *domain.User) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("postgres: RegisterUser: %w", err)
	}
	defer tx.Rollback()

	var exists bool
	if err := tx.QueryRowContext(ctx, `SELECT 1 FROM users WHERE user_name = $1`, user.UserName).Scan(&exists); err != nil {
		return fmt.Errorf("postgres: RegisterUser: %w", err)
	}
	if exists {
		return fmt.Errorf("postgres: RegisterUser: %w", domain.ErrUserNameExists)
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO users (user_name, password, created_at, updated_at) VALUES ($1, $2, $3, $4)`, user.UserName, user.Password, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("postgres: RegisterUser: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("postgres: RegisterUser: %w", err)
	}
	return nil
}

func (r *PSQLUsersRepository) User(ctx context.Context, userName string) (*domain.User, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("postgres: User: %w", err)
	}
	defer tx.Rollback()

	user := domain.User{UserName: userName}
	if err := tx.QueryRowContext(ctx, `SELECT user_id, password FROM users WHERE user_name = $1`, userName).Scan(&user.UserID, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("postgres: User: %w", domain.ErrUserNameNotExists)
		}
		return nil, fmt.Errorf("postgres: User: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("postgres: User: %w", err)
	}
	return &user, nil
}
