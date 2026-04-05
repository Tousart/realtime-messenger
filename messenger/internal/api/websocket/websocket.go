package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/tousart/messenger/internal/domain"
	"github.com/tousart/messenger/internal/dto"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:8080" || origin == "http://localhost:3000"
	},
}

type WebSocketRequest struct {
	Method  string          `json:"method"`
	Payload json.RawMessage `json:"payload"`
}

type Metadata struct {
	ctx    context.Context
	userID int
}

type WSMethod func(metadata *Metadata, conn *websocket.Conn, req *WebSocketRequest)

type WebSocketManager struct {
	messagesUC MessagesUsecase

	// // control reflection user an his current connection
	UserConnections map[int][]*websocket.Conn

	// control reflections chat-user
	ChatUsers map[int]map[int]int

	// control reflection user an his current chat
	UserChat map[int]int

	// for isolated access to maps
	mu *sync.RWMutex

	// methods
	WSMethods map[string]WSMethod

	// logger
	logger *slog.Logger
}

func NewWebSocketManager(messagesUC MessagesUsecase, logger *slog.Logger) *WebSocketManager {
	return &WebSocketManager{
		messagesUC:      messagesUC,
		UserConnections: make(map[int][]*websocket.Conn),
		ChatUsers:       make(map[int]map[int]int),
		UserChat:        make(map[int]int),
		mu:              &sync.RWMutex{},
		logger:          logger,
	}
}

func (ws *WebSocketManager) WithMethods() {
	ws.WSMethods = map[string]WSMethod{
		"send": WSMethod(ws.SendMessage),
		"join": WSMethod(ws.JoinToChat),
		// "leave": methods.MessengerMethod(LeaveChat),
	}
}

// type Client struct {
// 	conn     *websocket.Conn
// 	send     chan []byte
// 	metadata *Metadata
// 	ctx      context.Context
// 	cancel   context.CancelFunc
// }

func (ws *WebSocketManager) UpgradeConnectionForUser(w http.ResponseWriter, r *http.Request, responseHeader http.Header, payload *dto.UserPayload) error {
	const op = "websocket: UpgradeConnectionForUser:"

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ws.logger.Error("websocket: UpgradeConnectionForUser:", "err", err)
		return err
	}
	defer conn.Close()

	ws.connectUser(r.Context(), conn, payload)
	defer ws.disconnectUser(context.TODO(), conn, payload)

	metadata := Metadata{
		ctx:    r.Context(),
		userID: payload.ID,
	}

	for {
		_, wsRequest, err := conn.ReadMessage()
		if err != nil {
			ws.logger.Error(op, "err", err)
			break
		}

		ws.logger.Info("ws request:", "request", string(wsRequest))

		var req WebSocketRequest
		if err = json.Unmarshal(wsRequest, &req); err != nil {
			ws.logger.Info(op, fmt.Sprintf("invalid request from user %d:", payload.ID), err)
			continue
		}

		if method, ok := ws.WSMethods[req.Method]; ok {
			method(&metadata, conn, &req)
		} else {
			ws.logger.Info(op, "err", domain.ErrMethodNoTAllowed)
		}
	}

	return err
}

func (ws *WebSocketManager) SendMessageToUsersConnections(ctx context.Context, message dto.ConsumingMessage) error {
	// TODO: отправка сообщения в БД

	ws.mu.RLock()
	defer ws.mu.RUnlock()

	for userID := range ws.ChatUsers[message.ChatID] {
		for _, conn := range ws.UserConnections[userID] {
			if err := conn.WriteMessage(1, []byte(message.Text)); err != nil {
				ws.logger.Error("websocket: SendMessageToUsersConnections:", "err", err)
			}
		}
	}
	return nil
}

/*

	Подключение пользователя к WebSocket и отключение

*/

// connect user: add chats and connections
func (ws *WebSocketManager) connectUser(ctx context.Context, conn *websocket.Conn, payload *dto.UserPayload) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	// user - connection
	ws.UserConnections[payload.ID] = append(ws.UserConnections[payload.ID], conn)

	// chatID - user
	for _, chat := range payload.Chats {
		if _, ok := ws.ChatUsers[chat.ID]; !ok {
			if err := ws.messagesUC.SubscribeToChats(ctx, strconv.Itoa(chat.ID)); err != nil {
				ws.logger.Info("connectUser: failed subscribe user to chat:", "err", err)
				continue
			}
			ws.ChatUsers[chat.ID] = make(map[int]int)
		}
		ws.ChatUsers[chat.ID][payload.ID]++
	}
}

// disconnect user: remove chats and connections
func (ws *WebSocketManager) disconnectUser(ctx context.Context, conn *websocket.Conn, payload *dto.UserPayload) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	// user - connection
	connections := ws.UserConnections[payload.ID]

	if len(connections) == 1 {
		delete(ws.UserConnections, payload.ID)
	} else {
		for i, c := range connections {
			if c == conn {
				ws.UserConnections[payload.ID] = append(connections[:i], connections[i+1:]...)
				break
			}
		}
	}

	// chatID - user
	for _, chat := range payload.Chats {
		if ws.ChatUsers[chat.ID][payload.ID] == 1 {
			delete(ws.ChatUsers[chat.ID], payload.ID)
			if len(ws.ChatUsers[chat.ID]) == 1 {
				if err := ws.messagesUC.UnsubscribeFromChats(ctx, strconv.Itoa(chat.ID)); err != nil {
					ws.logger.Info("disconnectUser: failed unsubscribe user from chat:", "err", err)
					continue
				}
				delete(ws.ChatUsers, chat.ID)
			}
		} else {
			ws.ChatUsers[chat.ID][payload.ID]--
		}
	}
}
