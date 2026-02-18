package dto

type SendMessageWSRequest struct {
	UserID int    `json:"user_id"`
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

type ChatWSRequest struct {
	ChatID int `json:"chat_id"`
}

// ConsumingMessage - struct to receive messages (from broker) and send to users
type ConsumingMessage struct {
	UserID int    `json:"user_id"`
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}
