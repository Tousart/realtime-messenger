package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/tousart/messenger/internal/dto"
)

type ConsumeMessageResponse struct {
	UserID int    `json:"user_id"`
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == "http://localhost:8080"
	},
}

// Handle requests that becomes while websocket connection alive
func (ap *API) wsRequestHandler(ctx context.Context, conn *websocket.Conn, errChan chan<- error) {
	for {
		_, wsRequest, err := conn.ReadMessage()

		if err != nil {
			select {
			case <-ctx.Done():
			case errChan <- err:
			}
			return
		}

		log.Printf("ws request: %s\n", string(wsRequest))

		var req WebSocketRequest
		if err := json.Unmarshal(wsRequest, &req); err != nil {
			select {
			case <-ctx.Done():
			case errChan <- err:
			}
			return
		}

		if method, ok := ap.messengerMethods[req.Method]; ok {
			go method(req)
		} else {
			log.Printf("wsRequestHandler error: %v", errors.New("method not allowed"))
		}
	}
}

// TODO: реализовать нормальную логику консьюмера (ниже)
// добавить прослушивание очереди с repository вверх до api

// Consume messages and sending to users.
// Messages get from channel that exists
// only while websocket connection alive.
func (ap *API) send(msg *ConsumeMessageResponse) {
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

func (ap *API) consumeMessages(ctx context.Context, errChan chan<- error) {
	messagesQueue, err := ap.msgsHandlerService.MessagesQueue()
	if err != nil {
		select {
		case errChan <- fmt.Errorf("api: consumeMessages error: %s", err.Error()):
		case <-ctx.Done():
		}
		return
	}
	for {
		select {
		case msg, ok := <-messagesQueue:
			if !ok {
				log.Println("channel (queue) has been closed")
				return
			}
			msgBytes := msg.Body
			var message ConsumeMessageResponse
			if err := json.Unmarshal(msgBytes, &message); err != nil {
				log.Printf("api: consumeMessages error: %s\n", err.Error())
				continue
			}
			go ap.send(&message)
		case <-ctx.Done():
			return
		}
	}
}

// filling maps
// connect user: add chats and connections
// disconnect user: remove chats and connections

func (ap *API) connectUser(userID int, conn *websocket.Conn) {
	// TODO...
	// get users chats
	chats := Chats

	ap.mu.Lock()
	defer ap.mu.Unlock()

	// user - connection
	if ap.UserConnections[userID] == nil {
		ap.UserConnections[userID] = []*websocket.Conn{conn}
	} else {
		ap.UserConnections[userID] = append(ap.UserConnections[userID], conn)
	}

	// chatID - user
	for _, chatID := range chats {
		if _, ok := ap.ChatUsers[chatID]; !ok {
			ap.ChatUsers[chatID] = make(map[int]int)
			ap.msgsHandlerService.AddQueueToChat(context.Background(), dto.ChatWSRequest{ChatID: chatID})
		}
		ap.ChatUsers[chatID][userID]++
	}
}

func (ap *API) disconnectUser(userID int, conn *websocket.Conn) {
	// TODO...
	// get users chats
	chats := Chats

	ap.mu.Lock()
	defer ap.mu.Unlock()

	// user - connection
	connections := ap.UserConnections[userID]

	if len(connections) == 1 {
		delete(ap.UserConnections, userID)
	} else {
		for i, c := range connections {
			if c == conn {
				ap.UserConnections[userID] = append(connections[:i], connections[i+1:]...)
				break
			}
		}
	}

	// chatID - user
	for _, chatID := range chats {
		if ap.ChatUsers[chatID][userID] == 1 {
			delete(ap.ChatUsers[chatID], userID)
			if len(ap.ChatUsers[chatID]) == 0 {
				delete(ap.ChatUsers, chatID)
				ap.msgsHandlerService.RemoveQueueFromChat(context.Background(), dto.ChatWSRequest{ChatID: chatID})
			}
		} else {
			ap.ChatUsers[chatID][userID]--
		}
	}
}

// upgrade user connection to websocket

func (ap *API) messengerWebSocketConnectionHandler(w http.ResponseWriter, r *http.Request) {
	// TODO...
	// authorization
	// get userID
	userID := UserID

	// upgrade connection to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("messengerWebSocketConnectionHandler error: %s\n", err.Error())
		return
	}
	defer conn.Close()

	// ctx and errors channel for correct shutdown connection
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error)

	// connect and disconnect user
	// add users connection and chats to maps
	ap.connectUser(userID, conn)
	defer ap.disconnectUser(userID, conn)

	// handle requests on websocket connection
	go ap.wsRequestHandler(ctx, conn, errChan)

	// consume messages to this node
	go ap.consumeMessages(ctx, errChan)

	select {
	case <-ctx.Done():
	case err := <-errChan:
		log.Printf("messengerWebSocketConnectionHandler error: %s\n", err.Error())
	}
}
