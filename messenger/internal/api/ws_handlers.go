package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/tousart/messenger/internal/dto"
)

type WebSocketRequest struct {
	Method   string          `json:"method"`
	Metadata json.RawMessage `json:"metadata"`
	Payload  json.RawMessage `json:"payload"`
}

type Metadata struct {
	UserID int `json:"user_id"`
}

type MessengerMethod func(req WebSocketRequest)

func (ap *API) SendMessage(req WebSocketRequest) {
	var message dto.SendMessageWSRequest
	if err := json.Unmarshal(req.Payload, &message); err != nil {
		log.Printf("SendMessage error: %s\n", err.Error())
		return
	}

	var meta Metadata
	if err := json.Unmarshal(req.Metadata, &meta); err != nil {
		log.Printf("SendMessage error: %s\n", err.Error())
		return
	}
	message.UserID = meta.UserID

	if err := ap.msgsHandlerService.PublishMessageToQueues(context.TODO(), &message); err != nil {
		log.Printf("SendMessage error: %s\n", err.Error())
		return
	}
}

func (ap *API) JoinToChat(req WebSocketRequest) {
	var chat dto.ChatWSRequest
	if err := json.Unmarshal(req.Payload, &chat); err != nil {
		log.Printf("JoinToChat error: %s\n", err.Error())
		return
	}

	var meta Metadata
	if err := json.Unmarshal(req.Metadata, &meta); err != nil {
		log.Printf("SendMessage error: %s\n", err.Error())
		return
	}
	userID := meta.UserID

	ap.mu.RLock()
	if _, ok := ap.ChatUsers[chat.ChatID][userID]; !ok {
		log.Printf("JoinToChat error: %s\n", fmt.Errorf("user %d has not chat %d", userID, chat.ChatID).Error())
		return
	}
	ap.mu.RUnlock()

	ap.mu.Lock()
	ap.UserChat[userID] = chat.ChatID
	ap.mu.Unlock()
}
