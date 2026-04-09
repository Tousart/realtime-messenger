package dto

import "time"

type RegisterRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type User struct {
	ID        int64      `json:"user_id,string"`
	Name      string     `json:"user_name"`
	CreatedAt *time.Time `json:"created_at"`
}

type LoginRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type SessionID struct {
	SessionID string `json:"session_id"`
}

type SessionPayload struct {
	UserID   int64  `json:"user_id,string"`
	UserName string `json:"user_name"`
}

type ChatInfo struct {
	ID   int64  `json:"chat_id,string"`
	Name string `json:"chat_name"`
}

type UserPayload struct {
	ID    int64      `json:"user_id,string"`
	Name  string     `json:"user_name"`
	Chats []ChatInfo `json:"chats"`
}
