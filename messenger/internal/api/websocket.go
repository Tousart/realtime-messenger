package api

import (
	"context"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/tousart/messenger/internal/dto"
)

type WebSocketManager struct {
	UserConnections map[int][]*websocket.Conn

	// control reflections chat-user
	ChatUsers map[int]map[int]int

	// control reflection user an his current chat
	UserChat map[int]int

	// for isolated access to maps
	Mu *sync.RWMutex
}

func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		UserConnections: make(map[int][]*websocket.Conn),
		ChatUsers:       make(map[int]map[int]int),
		UserChat:        make(map[int]int),
		Mu:              &sync.RWMutex{},
	}
}

func (ws *WebSocketManager) SendMessageToUsersConnections(ctx context.Context, message dto.ConsumingMessage) error {
	ws.Mu.RLock()
	defer ws.Mu.RUnlock()

	for userID := range ws.ChatUsers[message.ChatID] {
		for _, conn := range ws.UserConnections[userID] {
			if err := conn.WriteMessage(1, []byte(message.Text)); err != nil {
				log.Printf("send error: %s\n", err.Error())
			}
		}
	}
	return nil
}
