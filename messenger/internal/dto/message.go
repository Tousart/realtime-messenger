package dto

import "time"

type SendMessageRequest struct {
	SenderID int64  `json:"sender_id,string"`
	ChatID   int64  `json:"chat_id,string"`
	Text     string `json:"text"`
}

type Message struct {
	ID        int64      `json:"message_id,string"`
	SenderID  int64      `json:"sender_id,string"`
	ChatID    int64      `json:"chat_id,string"`
	Text      string     `json:"text"`
	CreatedAt *time.Time `json:"created_at"`
}

type Messages struct {
	Messages []Message `json:"messages"`
}
