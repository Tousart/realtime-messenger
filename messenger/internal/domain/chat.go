package domain

import (
	"strings"
	"time"
)

type Chat struct {
	ID               int64
	Name             string
	ChatParticipants []ChatParticipant
	CreatedAt        *time.Time
}

type ChatParticipant struct {
	ID   int64
	Name string
	Role int
}

type ChatInfo struct {
	ID   int64
	Name string
}

func ValidateChatName(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrInvalidRequest
	}
	return nil
}
