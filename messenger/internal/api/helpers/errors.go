package helpers

import (
	"errors"
	"net/http"

	"github.com/tousart/messenger/internal/domain"
)

type errorStatus struct {
	msg    string
	status int
}

var errorToMessage = map[error]errorStatus{
	// Ошибки пользователя (400-404)
	domain.ErrUserNotFound:      {"user not found", http.StatusNotFound},
	domain.ErrUserAlreadyExists: {"user with this name already exists", http.StatusBadRequest},
	domain.ErrBadUserName:       {"invalid username format", http.StatusBadRequest},
	domain.ErrBadPassword:       {"password does not meet requirements", http.StatusBadRequest},
	domain.ErrIncorrectPassword: {"invalid password", http.StatusBadRequest},
	domain.ErrUnauthorized:      {"unauthorized", http.StatusUnauthorized},

	// Ошибки запроса
	domain.ErrMethodNoTAllowed: {"method not allowed", http.StatusMethodNotAllowed},
	domain.ErrInvalidRequest:   {"invalid request data", http.StatusBadRequest},
	domain.ErrEmptyChat:        {"chat is empty", http.StatusBadRequest},

	// Ошибки авторизации
	domain.ErrSessionIDNotExists: {"session expired or invalid", http.StatusUnauthorized},

	// Внутренние ошибки
	domain.ErrInternal: {"internal server error", http.StatusInternalServerError},
}

func MapError(err error) (string, int) {
	for targetErr, es := range errorToMessage {
		if errors.Is(err, targetErr) {
			return es.msg, es.status
		}
	}
	return "internal server error", http.StatusInternalServerError
}
