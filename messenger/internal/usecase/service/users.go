package service

import (
	"context"
	"fmt"

	"github.com/tousart/messenger/internal/domain"
	"github.com/tousart/messenger/internal/dto"
	"github.com/tousart/messenger/internal/repository"
	pkg "github.com/tousart/messenger/pkg/hashpassword"
)

type UsersService struct {
	usersRepo      repository.UsersRepository
	sessionsRepo   repository.SessionsRepository
	passwordHasher pkg.PasswordHasher
}

func NewUsersService(userRepo repository.UsersRepository, sessionsRepo repository.SessionsRepository, pswrdHasher pkg.PasswordHasher) *UsersService {
	return &UsersService{
		usersRepo:      userRepo,
		sessionsRepo:   sessionsRepo,
		passwordHasher: pswrdHasher,
	}
}

func (us *UsersService) RegisterUser(ctx context.Context, input dto.RegisterUserRequest) (*dto.RegisterUserResponse, error) {
	user, err := domain.NewUser(domain.WithUserName(input.UserName), domain.WithPassword(input.Password))
	if err != nil {
		return nil, fmt.Errorf("service: RegisterUser: %w", err)
	}

	hashedPassword, err := us.passwordHasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("service: RegisterUser: %w", err)
	}
	user.Password = hashedPassword

	userID, err := us.usersRepo.RegisterUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("service: RegisterUser: %w", err)
	}
	user.UserID = userID

	sessionID, err := us.sessionsRepo.GenerateSessionID(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("service: RegisterUser: %w", err)
	}
	return &dto.RegisterUserResponse{
		SessionID: sessionID,
	}, nil
}

func (us *UsersService) LoginUser(ctx context.Context, input dto.LoginUserRequest) (*dto.LoginUserResponse, error) {
	user, err := us.usersRepo.User(ctx, input.UserName)
	if err != nil {
		return nil, fmt.Errorf("service: LoginUser: %w", err)
	}

	if !us.passwordHasher.Compare(input.Password, user.Password) {
		return nil, fmt.Errorf("service: LoginUser: %w", domain.ErrIncorrectPassword)
	}

	sessionID, err := us.sessionsRepo.GenerateSessionID(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("service: RegisterUser: %w", err)
	}
	return &dto.LoginUserResponse{
		SessionID: sessionID,
	}, nil
}

func (us *UsersService) ValidateSessionID(ctx context.Context, sessionID string) (*dto.UserPayload, error) {
	user, err := us.sessionsRepo.SessionData(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("service: ValidateSessionID: %w", err)
	}
	return &dto.UserPayload{
		UserID:   user.UserID,
		UserName: user.UserName,
	}, nil
}
