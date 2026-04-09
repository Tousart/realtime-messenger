package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tousart/messenger/internal/domain"
	"github.com/tousart/messenger/internal/dto"
)

func (ws *WebSocketManager) SendMessageToConnections(message *dto.Message) {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	for userID := range ws.ChatUsers[message.ChatID] {
		for _, cw := range ws.UserConnections[userID] {
			ws.SendResponse(cw, TypeMessageNew, 0, message)
		}
	}
}

func (ws *WebSocketManager) SendMessage(metadata *Metadata, cw *ConnWriter, req *WebSocketRequest) {
	var sentMsg dto.SendMessageRequest
	err := json.Unmarshal(req.Payload, &sentMsg)
	if err != nil {
		ws.SendError(cw, domain.ErrInvalidRequest)
		return
	}

	sentMsg.SenderID = metadata.userID

	message, err := ws.messagesUC.SendMessage(metadata.ctx, &sentMsg)
	if err != nil {
		ws.SendError(cw, err)
		return
	}

	ws.SendResponse(cw, TypeMessageCreated, http.StatusCreated, message)
}

func (ws *WebSocketManager) JoinToChat(metadata *Metadata, cw *ConnWriter, req *WebSocketRequest) {
	var chat dto.JoinToChatRequest
	if err := json.Unmarshal(req.Payload, &chat); err != nil {
		ws.SendError(cw, domain.ErrInvalidRequest)
		return
	}

	ws.mu.RLock()
	_, ok := ws.ChatUsers[chat.ChatID][metadata.userID]
	ws.mu.RUnlock()

	if !ok {
		ws.SendError(cw, fmt.Errorf("user %d has not chat %d", metadata.userID, chat.ChatID))
		return
	}

	ws.mu.Lock()
	ws.UserChat[metadata.userID] = chat.ChatID
	ws.mu.Unlock()

	ws.SendResponse(cw, TypeOk, http.StatusOK, nil)
}

func (ws *WebSocketManager) CreateChat(metadata *Metadata, cw *ConnWriter, req *WebSocketRequest) {
	var createdChat dto.CreateChatRequest
	if err := json.Unmarshal(req.Payload, &createdChat); err != nil {
		ws.SendError(cw, domain.ErrInvalidRequest)
		return
	}

	if createdChat.ChatName == "" {
		ws.SendError(cw, fmt.Errorf("%w: chat name is required", domain.ErrInvalidRequest))
		return
	}

	createdChat.CreratorID = metadata.userID

	chat, err := ws.messagesUC.CreateChat(metadata.ctx, &createdChat)
	if err != nil {
		ws.SendError(cw, err)
		return
	}

	ws.SendResponse(cw, TypeChatCreated, http.StatusCreated, chat)
}
