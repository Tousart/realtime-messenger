package models

import "encoding/json"

type WebSocketMessageRequest struct {
	Method string          `json:"method"`
	Data   json.RawMessage `json:"data"`
}

type WSRequest struct {
	UserID int             `json:"user_id"`
	Data   json.RawMessage `json:"data"`
}

type Message struct {
	UserID int    `json:"user_id"`
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

type User struct {
	UserID int `json:"user_id"`
}

type Chat struct {
	ChatID int `json:"chat_id"`
}
