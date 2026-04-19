package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/tousart/messenger/internal/domain"
	"github.com/tousart/messenger/internal/dto"
)

type mockUsersRepo struct {
	createErr  error
	userErr    error
	returnUser *domain.User
}

func (m *mockUsersRepo) Create(_ context.Context, _ *domain.User) error {
	return m.createErr
}

func (m *mockUsersRepo) User(_ context.Context, _ string) (*domain.User, error) {
	if m.userErr != nil {
		return nil, m.userErr
	}
	if m.returnUser != nil {
		return m.returnUser, nil
	}
	return &domain.User{ID: 42, Name: "loba", Password: "hashed_password"}, nil
}

type mockSessionsRepo struct {
	genErr          error
	payloadErr      error
	returnSessionID string
}

func (m *mockSessionsRepo) GenerateSessionID(_ context.Context, _ *domain.SessionPayload) (string, error) {
	if m.genErr != nil {
		return "", m.genErr
	}

	if m.returnSessionID != "" {
		return m.returnSessionID, nil
	}

	return "test-session-id", nil
}

func (m *mockSessionsRepo) Payload(_ context.Context, _ string) (*domain.SessionPayload, error) {
	return nil, m.payloadErr
}

type mockHasher struct {
	hashErr error
}

func (m *mockHasher) Hash(_ string) (string, error) {
	if m.hashErr != nil {
		return "", m.hashErr
	}
	return "hashed_password", nil
}

func (m *mockHasher) Compare(_, _ string) bool { return true }

type mockIDGen struct{}

func (m *mockIDGen) GenerateID() int64 { return 42 }

/*
	──────────────────────────────────────────────────────────────
	Тесты
	──────────────────────────────────────────────────────────────
*/

// TestRegister_Success() проверяет успешную регистрацию с корректными данными.
func TestRegister_Success(t *testing.T) {
	uc := NewUsersService(
		&mockUsersRepo{},
		&mockSessionsRepo{},
		&mockHasher{},
		&mockIDGen{},
	)

	input := &dto.RegisterRequest{
		UserName: "alice",
		Password: "secret123",
	}

	user, err := uc.Register(context.Background(), input)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if user.Name != input.UserName {
		t.Errorf("expected name %q, got %q", input.UserName, user.Name)
	}
	if user.ID != 42 {
		t.Errorf("expected id 42, got %d", user.ID)
	}
	if user.CreatedAt == nil {
		t.Error("expected CreatedAt to be set")
	}
}

// TestRegister_InvalidInput() проверяет, что невалидные username/password
func TestRegister_InvalidInput(t *testing.T) {
	uc := NewUsersService(
		&mockUsersRepo{},
		&mockSessionsRepo{},
		&mockHasher{},
		&mockIDGen{},
	)

	type testcase struct {
		name  string
		input *dto.RegisterRequest
	}

	cases := []testcase{
		{"empty username", &dto.RegisterRequest{UserName: "", Password: "password"}},
		{"username with spaces", &dto.RegisterRequest{UserName: " loba", Password: "password"}},
		{"empty password", &dto.RegisterRequest{UserName: "loba", Password: ""}},
		{"password with spaces", &dto.RegisterRequest{UserName: "loba", Password: " password"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := uc.Register(context.Background(), tc.input)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !errors.Is(err, domain.ErrInvalidRequest) {
				t.Errorf("expected ErrInvalidRequest, got: %v", err)
			}
		})
	}
}

// TestRegister_RepoError() проверяет, что ошибка репозитория пробрасывается наружу
func TestRegister_RepoError(t *testing.T) {
	repoErr := domain.ErrUserAlreadyExists
	uc := NewUsersService(
		&mockUsersRepo{createErr: repoErr},
		&mockSessionsRepo{},
		&mockHasher{},
		&mockIDGen{},
	)

	input := &dto.RegisterRequest{
		UserName: "loba",
		Password: "password",
	}

	_, err := uc.Register(context.Background(), input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, repoErr) {
		t.Errorf("expected %v in error chain, got: %v", repoErr, err)
	}
}

// TestUser_Success() проверяет успешный логин: репозиторий возвращает пользователя
func TestLogin_Success(t *testing.T) {
	user := &domain.User{
		ID:       42,
		Name:     "vika",
		Password: "secret",
	}

	sessionID := "token"

	uc := NewUsersService(
		&mockUsersRepo{
			returnUser: user,
		},
		&mockSessionsRepo{
			returnSessionID: sessionID,
		},
		&mockHasher{},
		&mockIDGen{},
	)

	input := &dto.LoginRequest{
		UserName: "vika",
		Password: "secret",
	}

	session, err := uc.Login(context.Background(), input)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if session.SessionID == "" {
		t.Error("expected non-empty session ID")
	}

	if session.SessionID != sessionID {
		t.Error("invalid session id")
	}
}

// TestRegister_RepoError() проверяет, что ошибка репозитория пробрасывается наружу
func TestLogin_RepoError(t *testing.T) {
	repoErr := domain.ErrUserNotFound
	uc := NewUsersService(
		&mockUsersRepo{userErr: repoErr},
		&mockSessionsRepo{},
		&mockHasher{},
		&mockIDGen{},
	)

	input := &dto.LoginRequest{
		UserName: "vika",
		Password: "secret",
	}

	_, err := uc.Login(context.Background(), input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, repoErr) {
		t.Errorf("expected %v in error chain, got: %v", repoErr, err)
	}
}
