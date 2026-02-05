package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/tousart/messenger/internal/models"
)

func (ap *API) SendMessage(req models.WSRequest) {
	var message models.Message
	if err := json.Unmarshal(req.Data, &message); err != nil {
		log.Printf("SendMessage error: %s\n", err.Error())
		return
	}

	message.UserID = req.UserID

	if err := ap.publisherService.PublishMessage(context.TODO(), message); err != nil {
		log.Printf("SendMessage error: %s\n", err.Error())
		return
	}
}

func (ap *API) JoinToChat(req models.WSRequest) {
	var chat models.Chat
	if err := json.Unmarshal(req.Data, &chat); err != nil {
		log.Printf("JoinToChat error: %s\n", err.Error())
		return
	}

	userID := req.UserID

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
