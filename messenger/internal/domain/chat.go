package domain

import "strings"

type Chat struct {
	ID               int
	Name             string
	ChatParticipants []ChatParticipant
}

type ChatParticipant struct {
	UserID   int
	UserName string
	Role     int
}

type ChatInfo struct {
	ID   int
	Name string
}

func ValidateChatName(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrInvalidRequest
	}
	return nil
}
