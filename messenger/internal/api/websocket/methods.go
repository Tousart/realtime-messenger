package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/tousart/messenger/internal/api/helpers"
	"github.com/tousart/messenger/internal/domain"
	"github.com/tousart/messenger/internal/dto"
)

func (ws *WebSocketManager) SendMessage(metadata *Metadata, conn *websocket.Conn, req *WebSocketRequest) {
	var sentMsg dto.SendMessageRequest
	err := json.Unmarshal(req.Payload, &sentMsg)
	if err != nil {
		ws.SendError(conn, domain.ErrInvalidRequest)
		return
	}

	sentMsg.SenderID = metadata.userID

	message, err := ws.messagesUC.SendMessage(metadata.ctx, &sentMsg)
	if err != nil {
		ws.SendError(conn, err)
		return
	}

	ws.SendResponse(conn, http.StatusCreated, message)
}

func (ws *WebSocketManager) JoinToChat(metadata *Metadata, conn *websocket.Conn, req *WebSocketRequest) {
	var chat dto.JoinToChatRequest
	if err := json.Unmarshal(req.Payload, &chat); err != nil {
		ws.SendError(conn, domain.ErrInvalidRequest)
		return
	}

	ws.mu.RLock()
	if _, ok := ws.ChatUsers[chat.ChatID][metadata.userID]; !ok {
		ws.SendError(conn, fmt.Errorf("user %d has not chat %d", metadata.userID, chat.ChatID))
		return
	}
	ws.mu.RUnlock()

	ws.mu.Lock()
	ws.UserChat[metadata.userID] = chat.ChatID
	ws.mu.Unlock()

	ws.SendResponse(conn, http.StatusOK, nil)
}

func (ws *WebSocketManager) CreateChat(metadata *Metadata, conn *websocket.Conn, req *WebSocketRequest) {
	var createdChat dto.CreateChatRequest
	if err := json.Unmarshal(req.Payload, &createdChat); err != nil {
		ws.SendError(conn, domain.ErrInvalidRequest)
		return
	}

	if createdChat.ChatName == "" {
		ws.SendError(conn, fmt.Errorf("%w: chat name is required", domain.ErrInvalidRequest))
		return
	}

	createdChat.CreratorID = metadata.userID

	chat, err := ws.messagesUC.CreateChat(metadata.ctx, &createdChat)
	if err != nil {
		ws.SendError(conn, err)
		return
	}

	ws.SendResponse(conn, http.StatusCreated, chat)
}

/*

	Ответ

*/

type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

func (ws *WebSocketManager) SendError(conn *websocket.Conn, err error) {
	if errors.Is(err, domain.ErrInternal) {
		ws.logger.Error("ERROR:", "err", err)
	} else {
		ws.logger.Info("error:", "err", err)
	}
	msg, status := helpers.MapError(err)
	ws.Send(conn, ErrorResponse{
		Status: status,
		Error:  msg,
	})
}

type Response struct {
	Status int `json:"status"`
	Body   any `json:"body,omitempty"`
}

func (ws *WebSocketManager) SendResponse(conn *websocket.Conn, status int, body any) {
	ws.Send(conn, Response{
		Status: status,
		Body:   body,
	})
}

func (ws *WebSocketManager) Send(conn *websocket.Conn, response any) {
	msg, err := json.Marshal(response)
	if err != nil {
		ws.logger.Error("websocket: Send:", "err", err)
	}

	if err = conn.WriteMessage(1, msg); err != nil {
		ws.logger.Error("websocket: Send:", "err", err)
	}
}
