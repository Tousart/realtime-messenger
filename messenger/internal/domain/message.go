package domain

import (
	"strings"
	"time"
)

type Message struct {
	ID        int64
	SenderID  int64
	ChatID    int64
	Text      string
	CreatedAt *time.Time
}

func IsValidMessageText(text string) (string, error) {
	t := strings.TrimSpace(text)
	if len(t) == 0 {
		return "", ErrInvalidRequest
	}
	return t, nil
}
