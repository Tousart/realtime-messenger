package api

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/tousart/messenger/internal/dto"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == "http://localhost:8080"
	},
}

/*

	Обработчик-маршрутизатор методов внутри WebSocket-соединения

*/

// Handle requests that becomes while websocket connection alive
func (ap *API) wsRequestHandler(ctx context.Context, conn *websocket.Conn, metadata *Metadata, errChan chan<- error) {
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

		req := &WebSocketRequest{}
		if err := json.Unmarshal(wsRequest, req); err != nil {
			select {
			case <-ctx.Done():
			case errChan <- err:
			}
			return
		}

		if method, ok := ap.messengerMethods[req.Method]; ok {
			go method(metadata, req)
		} else {
			log.Printf("wsRequestHandler error: %v", errors.New("method not allowed"))
		}
	}
}

/*

	Подключение пользователя к WebSocket и отключение соответственно

*/

// connect user: add chats and connections
func (ap *API) connectUser(userID int, conn *websocket.Conn) {
	// TODO...
	// get users chats
	chats := Chats

	ap.wsManager.Mu.Lock()
	defer ap.wsManager.Mu.Unlock()

	// user - connection
	ap.wsManager.UserConnections[userID] = append(ap.wsManager.UserConnections[userID], conn)

	// chatID - user
	for _, chatID := range chats {
		if _, ok := ap.wsManager.ChatUsers[chatID]; !ok {
			if err := ap.msgsHandlerService.AddQueueToChat(context.Background(), dto.ChatWSRequest{ChatID: chatID}); err != nil {
				log.Printf("connectUser: failed add queue to chat: %v\n", err)
				continue
			}
			ap.wsManager.ChatUsers[chatID] = make(map[int]int)
		}
		ap.wsManager.ChatUsers[chatID][userID]++
	}
}

// disconnect user: remove chats and connections
func (ap *API) disconnectUser(userID int, conn *websocket.Conn) {
	// TODO...
	// get users chats
	chats := Chats

	ap.wsManager.Mu.Lock()
	defer ap.wsManager.Mu.Unlock()

	// user - connection
	connections := ap.wsManager.UserConnections[userID]

	if len(connections) == 1 {
		delete(ap.wsManager.UserConnections, userID)
	} else {
		for i, c := range connections {
			if c == conn {
				ap.wsManager.UserConnections[userID] = append(connections[:i], connections[i+1:]...)
				break
			}
		}
	}

	// chatID - user
	for _, chatID := range chats {
		if ap.wsManager.ChatUsers[chatID][userID] == 1 {
			if len(ap.wsManager.ChatUsers[chatID]) == 1 {
				if err := ap.msgsHandlerService.RemoveQueueFromChat(context.Background(), dto.ChatWSRequest{ChatID: chatID}); err != nil {
					log.Printf("disconnectUser: failed remove queue from chat: %v\n", err)
					continue
				}
				delete(ap.wsManager.ChatUsers, chatID)
			}
			delete(ap.wsManager.ChatUsers[chatID], userID)
		} else {
			ap.wsManager.ChatUsers[chatID][userID]--
		}
	}
}

/*

	Обработчик улучшения соединения до WebSocket

*/

func (ap *API) messengerWebSocketConnectionHandler(w http.ResponseWriter, r *http.Request) {
	// TODO...
	// authorization
	// get userID
	meta := Metadata{
		UserID: UserID,
	}
	// userID := UserID

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
	ap.connectUser(meta.UserID, conn)
	defer ap.disconnectUser(meta.UserID, conn)

	// handle requests on websocket connection
	go ap.wsRequestHandler(ctx, conn, &meta, errChan)

	// consume messages to this node
	// go ap.consumeMessages(ctx, errChan)

	select {
	case <-ctx.Done():
	case err := <-errChan:
		log.Printf("messengerWebSocketConnectionHandler error: %s\n", err.Error())
	}
}
