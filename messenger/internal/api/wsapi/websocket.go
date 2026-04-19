package wsapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/tousart/messenger/internal/domain"
	"github.com/tousart/messenger/internal/dto"
	"github.com/tousart/messenger/internal/middleware/wsmw"
	"github.com/tousart/messenger/pkg/types/wstypes"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:8080" || origin == "http://localhost:3000"
	},
}

type WebSocketManager struct {
	messagesUC MessagesUsecase

	// // control reflection user an his current connection
	UserConnections map[int64][]*wstypes.ConnWriter

	// control reflections chat-user
	ChatUsers map[int64]map[int64]int

	// for isolated access to maps
	mu *sync.RWMutex

	// methods
	WSMethods map[string]wstypes.Method

	// logger
	logger *slog.Logger
}

func NewWebSocketManager(messagesUC MessagesUsecase, logger *slog.Logger) *WebSocketManager {
	return &WebSocketManager{
		messagesUC:      messagesUC,
		UserConnections: make(map[int64][]*wstypes.ConnWriter),
		ChatUsers:       make(map[int64]map[int64]int),
		mu:              &sync.RWMutex{},
		WSMethods:       make(map[string]wstypes.Method),
		logger:          logger,
	}
}

func (ws *WebSocketManager) addWSMethod(name string, method wstypes.Method) {
	ws.WSMethods[name] = method
}

func (ws *WebSocketManager) useMiddleware(mw func(wstypes.Method) wstypes.Method) {
	for name, method := range ws.WSMethods {
		ws.WSMethods[name] = mw(method)
	}
}

func (ws *WebSocketManager) WithMethods() {
	ws.addWSMethod("send", ws.SendMessage)
	ws.addWSMethod("join", ws.JoinToChat)
	ws.addWSMethod("create", ws.CreateChat)

	ws.useMiddleware(wsmw.Logging(ws.logger))
}

func (ws *WebSocketManager) UpgradeConnectionForUser(w http.ResponseWriter, r *http.Request, responseHeader http.Header, payload *dto.UserPayload) error {
	const op = "websocket: UpgradeConnectionForUser:"

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ws.logger.Error(op, "err", err)
		return err
	}
	cw := wstypes.NewConnWriter(conn)
	defer cw.Close()

	ws.connectUser(r.Context(), cw, payload)
	defer ws.disconnectUser(context.TODO(), cw, payload)

	if err := wstypes.Send(cw, payload); err != nil {
		ws.logger.Error(op, "err", err)
		return err
	}

	metadata := wstypes.Metadata{
		Ctx:    r.Context(),
		UserID: payload.ID,
	}

	for {
		_, wsRequest, err := conn.ReadMessage()
		if err != nil {
			ws.logger.Error(op, "err", err)
			break
		}

		ws.logger.Info("ws request:", "request", string(wsRequest))

		var req wstypes.Request
		if err = json.Unmarshal(wsRequest, &req); err != nil {
			ws.logger.Info(op, fmt.Sprintf("invalid request from user %d:", payload.ID), err)
			continue
		}

		if method, ok := ws.WSMethods[req.Method]; ok {
			method(&metadata, cw, &req)
		} else {
			ws.logger.Info(op, "err", domain.ErrMethodNoTAllowed)
		}
	}

	return err
}

/*

	Подключение пользователя к WebSocket и отключение

*/

// connect user: add chats and connections
func (ws *WebSocketManager) connectUser(ctx context.Context, cw *wstypes.ConnWriter, payload *dto.UserPayload) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	// user - connection
	ws.UserConnections[payload.ID] = append(ws.UserConnections[payload.ID], cw)

	// chatID - user
	for _, chat := range payload.Chats {
		if _, ok := ws.ChatUsers[chat.ID]; !ok {
			if err := ws.messagesUC.SubscribeToChats(ctx, chat.ID); err != nil {
				ws.logger.Info("connectUser: failed subscribe user to chat:", "err", err)
				continue
			}
			ws.ChatUsers[chat.ID] = make(map[int64]int)
		}
		ws.ChatUsers[chat.ID][payload.ID]++
	}
}

// disconnect user: remove chats and connections
func (ws *WebSocketManager) disconnectUser(ctx context.Context, cw *wstypes.ConnWriter, payload *dto.UserPayload) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	// user - connection
	connections := ws.UserConnections[payload.ID]

	if len(connections) == 1 {
		delete(ws.UserConnections, payload.ID)
	} else {
		for i, c := range connections {
			if c == cw {
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
				if err := ws.messagesUC.UnsubscribeFromChats(ctx, chat.ID); err != nil {
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
