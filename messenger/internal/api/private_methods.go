package api

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
)

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
			// TODO: пул горутин, которые будут параллельно обрабатывать запросы пользователя
			// go method(metadata, req)
			method(metadata, req)
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
			if err := ap.msgsHandlerService.SubscribeToChats(context.Background(), strconv.Itoa(chatID)); err != nil {
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
				if err := ap.msgsHandlerService.UnsubscribeFromChats(context.Background(), strconv.Itoa(chatID)); err != nil {
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
