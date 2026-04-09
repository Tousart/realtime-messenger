package dto

import "time"

type CreateChatRequest struct {
	ChatName         string            `json:"chat_name"`
	CreratorID       int64             `json:"creator_id,string"`
	ChatParticipants []ChatParticipant `json:"chat_participants"`
}

type Chat struct {
	ID               int64             `json:"chat_id,string"`
	Name             string            `json:"chat_name"`
	ChatParticipants []ChatParticipant `json:"chat_participants"`
	CreatedAt        *time.Time        `json:"created_at"`
}

type ChatParticipant struct {
	ID   int64   `json:"participant_id,string"`
	Name *string `json:"participant_name"`
	Role *int    `json:"role"`
}

type JoinToChatRequest struct {
	ChatID int64 `json:"chat_id,string"`
}
