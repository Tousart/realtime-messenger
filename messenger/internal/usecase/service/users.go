package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/tousart/messenger/internal/domain"
	"github.com/tousart/messenger/internal/repository"
	"github.com/tousart/messenger/pkg"
)

const MAX_PASSWORD_LENGTH = 72

type UsersService struct {
	usersRepo      repository.UsersRepository
	passwordHasher pkg.PasswordHasher
}

func NewUsersService(userRepo repository.UsersRepository, pswrdHasher pkg.PasswordHasher) *UsersService {
	return &UsersService{
		usersRepo:      userRepo,
		passwordHasher: pswrdHasher,
	}
}

func (us *UsersService) RegisterUser(ctx context.Context, data *domain.RegisterRequest) error {
	if len(strings.TrimSpace(data.Password)) > MAX_PASSWORD_LENGTH {
		return fmt.Errorf("service: RegisterUser: %w", domain.ErrBadPassword)
	}

	hashedPassword, err := us.passwordHasher.Hash(data.Password)
	if err != nil {
		return fmt.Errorf("service: RegisterUser: %w", err)
	}
	user := domain.User{
		UserName: data.UserName,
		Password: hashedPassword,
	}

	if err := us.usersRepo.RegisterUser(ctx, &user); err != nil {
		return fmt.Errorf("service: RegisterUser: %w", err)
	}
	return nil
}

func (us *UsersService) LoginUser(ctx context.Context, data *domain.LoginRequest) error {
	user, err := us.usersRepo.User(ctx, data.UserName)
	if err != nil {
		return fmt.Errorf("service: LoginUser: %w", err)
	}

	if !us.passwordHasher.Compare(data.Password, user.Password) {
		return fmt.Errorf("service: LoginUser: %w", domain.ErrIncorrectPassword)
	}
	return nil
}
