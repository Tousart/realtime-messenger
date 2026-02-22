package domain

import "errors"

var (
	ErrMethodNoTAllowed error = errors.New("method not allowed")

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
