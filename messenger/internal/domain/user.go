package domain

import (
	"strings"
	"time"
)

type User struct {
	ID        int64
	Name      string
	Password  string
	CreatedAt *time.Time
}

func IsValidUserName(name string) error {
	if len(name) == 0 || len(strings.TrimSpace(name)) != len(name) {
		return ErrBadUserName
	}
	return nil
}

func IsValidUserPassword(password string) error {
	if len(password) == 0 || len(strings.TrimSpace(password)) != len(password) {
		return ErrBadPassword
	}
	return nil
}
