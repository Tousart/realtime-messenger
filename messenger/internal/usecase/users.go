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
	idGen          IDGenerator
}

func NewUsersService(userRepo UsersRepository, sessionsRepo SessionsRepository, pswrdHasher PasswordHasher, idGen IDGenerator) *UsersUsecase {
	return &UsersUsecase{
		usersRepo:      userRepo,
		sessionsRepo:   sessionsRepo,
		passwordHasher: pswrdHasher,
		idGen:          idGen,
	}
}

func (uc *UsersUsecase) Register(ctx context.Context, input *dto.RegisterRequest) (*dto.User, error) {
	const op = "usecase: Register:"

	err := domain.IsValidUserName(input.UserName)
	if err != nil {
		return nil, fmt.Errorf("%s %w: %w", op, domain.ErrInvalidRequest, err)
	}

	if err = domain.IsValidUserPassword(input.Password); err != nil {
		return nil, fmt.Errorf("%s %w: %w", op, domain.ErrInvalidRequest, err)
	}

	hashedPassword, err := uc.passwordHasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	createdAt := timeNowUTC()

	user := &domain.User{
		ID:        uc.idGen.GenerateID(),
		Name:      input.UserName,
		Password:  hashedPassword,
		CreatedAt: &createdAt,
	}

	if err = uc.usersRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	return &dto.User{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (uc *UsersUsecase) Login(ctx context.Context, input *dto.LoginRequest) (*dto.SessionID, error) {
	const op = "usecase: Login:"

	user, err := uc.usersRepo.User(ctx, input.UserName)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	if !uc.passwordHasher.Compare(user.Password, input.Password) {
		return nil, fmt.Errorf("%s %w", op, domain.ErrIncorrectPassword)
	}

	sessionPayload := &domain.SessionPayload{
		UserID:   user.ID,
		UserName: user.Name,
	}

	sessionID, err := uc.sessionsRepo.GenerateSessionID(ctx, sessionPayload)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	return &dto.SessionID{
		SessionID: sessionID,
	}, nil
}

func (uc *UsersUsecase) ValidateSessionID(ctx context.Context, sessionID string) (*dto.SessionPayload, error) {
	payload, err := uc.sessionsRepo.Payload(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("usecase: ValidateSessionID: %w", err)
	}

	return &dto.SessionPayload{
		UserID:   payload.UserID,
		UserName: payload.UserName,
	}, nil
}
