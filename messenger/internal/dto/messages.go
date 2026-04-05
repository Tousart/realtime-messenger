package dto

// SendMessage
type SendMessageRequest struct {
	UserID int    `json:"user_id"`
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

// JoinToChat
type JoinToChatRequest struct {
	ChatID int `json:"chat_id"`
}

// CreateChat
type CreateChatRequest struct {
	ChatName         string            `json:"chat_name"`
	CreratorID       int               `json:"creator_id"`
	ChatParticipants []ChatParticipant `json:"chat_participants"`
}

type ChatParticipant struct {
	UserName string `json:"user_name"`
}

type CreateChatResponse struct {
	ChatID           int                       `json:"chat_id"`
	ChatName         string                    `json:"chat_name"`
	ChatParticipants []ChatParticipantResponse `json:"chat_participants"`
}

type ChatParticipantResponse struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
	Role     int    `json:"role"`
}

type ConsumingMessage struct {
	UserID int    `json:"user_id"`
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

type Chat struct {
	ID   int    `json:"chat_id"`
	Name string `json:"chat_name"`
}

type UserPayload struct {
	ID    int    `json:"user_id"`
	Name  string `json:"user_name"`
	Chats []Chat `json:"chats"`
}
