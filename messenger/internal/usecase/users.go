package usecase

import (
	"context"
	"fmt"

	"github.com/tousart/messenger/internal/domain"
	"github.com/tousart/messenger/internal/dto"
)

type UsersUsecase struct {
	usersRepo      UsersRepository
	sessionsRepo   SessionsRepository
	passwordHasher PasswordHasher
}

func NewUsersService(userRepo UsersRepository, sessionsRepo SessionsRepository, pswrdHasher PasswordHasher) *UsersUsecase {
	return &UsersUsecase{
		usersRepo:      userRepo,
		sessionsRepo:   sessionsRepo,
		passwordHasher: pswrdHasher,
	}
}

func (u *UsersUsecase) RegisterUser(ctx context.Context, input dto.RegisterUserRequest) (*dto.RegisterUserResponse, error) {
	user, err := domain.NewUser(domain.WithUserName(input.UserName), domain.WithPassword(input.Password))
	if err != nil {
		return nil, fmt.Errorf("service: RegisterUser: %w", err)
	}

	hashedPassword, err := u.passwordHasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("service: RegisterUser: %w", err)
	}
	user.Password = hashedPassword

	userID, err := u.usersRepo.RegisterUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("service: RegisterUser: %w", err)
	}
	user.UserID = userID

	sessionID, err := u.sessionsRepo.GenerateSessionID(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("service: RegisterUser: %w", err)
	}
	return &dto.RegisterUserResponse{
		SessionID: sessionID,
	}, nil
}

func (u *UsersUsecase) LoginUser(ctx context.Context, input dto.LoginUserRequest) (*dto.LoginUserResponse, error) {
	user, err := u.usersRepo.User(ctx, input.UserName)
	if err != nil {
		return nil, fmt.Errorf("service: LoginUser: %w", err)
	}

	if !u.passwordHasher.Compare(user.Password, input.Password) {
		return nil, fmt.Errorf("service: LoginUser: %w", domain.ErrIncorrectPassword)
	}

	sessionID, err := u.sessionsRepo.GenerateSessionID(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("service: RegisterUser: %w", err)
	}
	return &dto.LoginUserResponse{
		SessionID: sessionID,
	}, nil
}

func (u *UsersUsecase) ValidateSessionID(ctx context.Context, sessionID string) (*dto.SessionPayload, error) {
	user, err := u.sessionsRepo.SessionIDPayload(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("service: ValidateSessionID: %w", err)
	}
	return &dto.SessionPayload{
		UserID:   user.UserID,
		UserName: user.UserName,
	}, nil
}
