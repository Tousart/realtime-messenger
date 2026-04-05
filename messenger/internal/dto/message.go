package dto

import "time"

type SendMessageRequest struct {
	SenderID int64  `json:"sender_id"`
	ChatID   int64  `json:"chat_id"`
	Text     string `json:"text"`
}

type Message struct {
	SenderID  int64      `json:"sender_id"`
	ChatID    int64      `json:"chat_id"`
	Text      string     `json:"text"`
	CreatedAt *time.Time `created_at:"json"`
}
