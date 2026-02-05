package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/tousart/messenger/internal/models"
)

func (ap *API) receiveMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var msg models.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	go ap.send(&msg)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (ap *API) send(msg *models.Message) {
	ap.mu.RLock()
	defer ap.mu.RUnlock()

	for userID := range ap.ChatUsers[msg.ChatID] {
		for _, conn := range ap.UserConnections[userID] {
			if err := conn.WriteMessage(1, []byte(msg.Text)); err != nil {
				log.Printf("send error: %s\n", err.Error())
			}
		}
	}
}
