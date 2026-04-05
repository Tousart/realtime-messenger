package domain

import "errors"

var (
	ErrMethodNoTAllowed error = errors.New("method not allowed")
	ErrInvalidRequest   error = errors.New("invalid request")
	ErrInternal         error = errors.New("internal error")
	ErrUnauthorized     error = errors.New("unauthorized")

	// ErrBadChatName error = errors.New("bad chat name")
	// ErrBadChatID   error = errors.New("bad chat id")
	ErrEmptyChat error = errors.New("chat is empty")

	ErrSessionIDNotExists error = errors.New("session id not exists")

	// USERS ERRORS
	// username
	ErrUserExists   error = errors.New("user exists")
	ErrUserNotFound error = errors.New("user not found")
	ErrBadUserName  error = errors.New("bad username")
	// password
	ErrBadPassword       error = errors.New("bad password")
	ErrIncorrectPassword error = errors.New("incorrect password")
)
