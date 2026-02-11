package domain

import "errors"

var (
	// USERS ERRORS

	// username
	ErrUserNameExists    error = errors.New("username exists")
	ErrUserNameNotExists error = errors.New("username not exists")
	// password
	ErrBadPassword       error = errors.New("bad password")
	ErrIncorrectPassword error = errors.New("incorrect password")
)
