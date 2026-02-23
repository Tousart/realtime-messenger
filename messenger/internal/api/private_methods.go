package api

import (
	"context"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
)

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
				log.Printf("connectUser: failed subscribe user to chat: %v\n", err)
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
					log.Printf("disconnectUser: failed unsubscribe user from chat: %v\n", err)
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
