package dto

import "time"

type SendMessageRequest struct {
	SenderID int64  `json:"sender_id,string"`
	ChatID   int64  `json:"chat_id,string"`
	Text     string `json:"text"`
}

type Message struct {
	SenderID  int64      `json:"sender_id,string"`
	ChatID    int64      `json:"chat_id,string"`
	Text      string     `json:"text"`
	CreatedAt *time.Time `json:"created_at"`
}
