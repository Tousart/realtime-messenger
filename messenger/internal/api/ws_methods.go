package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/tousart/messenger/internal/dto"
)

type WebSocketRequest struct {
	Method  string          `json:"method"`
	Payload json.RawMessage `json:"payload"`
}

type Metadata struct {
	UserID int `json:"user_id"`
}

type MessengerMethod func(metadata *Metadata, req *WebSocketRequest)

func (ap *API) SendMessage(metadata *Metadata, req *WebSocketRequest) {
	var message dto.SendMessageWSRequest
	if err := json.Unmarshal(req.Payload, &message); err != nil {
		log.Printf("SendMessage error: %s\n", err.Error())
		return
	}

	message.UserID = metadata.UserID

	if err := ap.msgsHandlerService.PublishMessageToQueues(context.TODO(), message); err != nil {
		log.Printf("SendMessage error: %s\n", err.Error())
		return
	}
}

func (ap *API) JoinToChat(metadata *Metadata, req *WebSocketRequest) {
	var chat dto.ChatWSRequest
	if err := json.Unmarshal(req.Payload, &chat); err != nil {
		log.Printf("JoinToChat error: %s\n", err.Error())
		return
	}

	ap.wsManager.Mu.RLock()
	if _, ok := ap.wsManager.ChatUsers[chat.ChatID][metadata.UserID]; !ok {
		log.Printf("JoinToChat error: %s\n", fmt.Errorf("user %d has not chat %d", metadata.UserID, chat.ChatID).Error())
		return
	}
	ap.wsManager.Mu.RUnlock()

	ap.wsManager.Mu.Lock()
	ap.wsManager.UserChat[metadata.UserID] = chat.ChatID
	ap.wsManager.Mu.Unlock()
}
