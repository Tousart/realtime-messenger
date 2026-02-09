package models

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type WebSocketMessageRequest struct {
	Method string          `json:"method"`
	Data   json.RawMessage `json:"data"`
}

type WSRequest struct {
	UserID int             `json:"user_id"`
	Data   json.RawMessage `json:"data"`
}

type MessagesQueue <-chan amqp.Delivery

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
