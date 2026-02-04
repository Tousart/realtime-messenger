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

	log.Printf("message text: %s\n", message.Text)

	if err := ap.publisherService.PublishMessage(context.TODO(), message); err != nil {
		log.Printf("SendMessage error: %s\n", err.Error())
		return
	}

	// ap.send(&message)
}

// func (ap *API) send(msg *models.Message) {
// 	ap.mu.RLock()
// 	defer ap.mu.RUnlock()

// 	for userID := range ap.ChatUsers[msg.ChatID] {
// 		for _, conn := range ap.UserConnections[userID] {
// 			if err := conn.WriteMessage(1, []byte(msg.Text)); err != nil {
// 				log.Printf("send error: %s\n", err.Error())
// 			}
// 		}
// 	}
// }

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
