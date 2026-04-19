package wsapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tousart/messenger/internal/api/helpers"
	"github.com/tousart/messenger/internal/domain"
	"github.com/tousart/messenger/internal/dto"
	"github.com/tousart/messenger/pkg/types/wstypes"
)

func (ws *WebSocketManager) SendMessageToConnections(message *dto.Message) {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	for userID := range ws.ChatUsers[message.ChatID] {
		for _, cw := range ws.UserConnections[userID] {
			ws.SendResponse(cw, wstypes.TypeMessageNew, 0, message)
		}
	}
}

func (ws *WebSocketManager) SendMessage(metadata *wstypes.Metadata, cw *wstypes.ConnWriter, req *wstypes.Request) {
	var sentMsg dto.SendMessageRequest
	err := json.Unmarshal(req.Payload, &sentMsg)
	if err != nil {
		ws.SendError(cw, domain.ErrInvalidRequest)
		return
	}

	sentMsg.SenderID = metadata.UserID

	message, err := ws.messagesUC.SendMessage(metadata.Ctx, &sentMsg)
	if err != nil {
		ws.SendError(cw, err)
		return
	}

	ws.SendResponse(cw, wstypes.TypeMessageCreated, http.StatusCreated, message)
}

func (ws *WebSocketManager) JoinToChat(metadata *wstypes.Metadata, cw *wstypes.ConnWriter, req *wstypes.Request) {
	var chat dto.JoinToChatRequest
	if err := json.Unmarshal(req.Payload, &chat); err != nil {
		ws.SendError(cw, domain.ErrInvalidRequest)
		return
	}

	ws.mu.RLock()
	_, ok := ws.ChatUsers[chat.ChatID][metadata.UserID]
	ws.mu.RUnlock()

	if !ok {
		ws.SendError(cw, fmt.Errorf("user %d has not chat %d", metadata.UserID, chat.ChatID))
		return
	}

	messages, err := ws.messagesUC.Messages(metadata.Ctx, chat.ChatID)
	if err != nil {
		ws.SendError(cw, err)
	}

	ws.SendResponse(cw, wstypes.TypeOk, http.StatusOK, messages)
}

func (ws *WebSocketManager) CreateChat(metadata *wstypes.Metadata, cw *wstypes.ConnWriter, req *wstypes.Request) {
	var createdChat dto.CreateChatRequest
	if err := json.Unmarshal(req.Payload, &createdChat); err != nil {
		ws.SendError(cw, domain.ErrInvalidRequest)
		return
	}

	if createdChat.ChatName == "" {
		ws.SendError(cw, fmt.Errorf("%w: chat name is required", domain.ErrInvalidRequest))
		return
	}

	createdChat.CreratorID = metadata.UserID

	chat, err := ws.messagesUC.CreateChat(metadata.Ctx, &createdChat)
	if err != nil {
		ws.SendError(cw, err)
		return
	}

	ws.mu.Lock()
	ws.ChatUsers[chat.ID] = make(map[int64]int)
	for _, participant := range chat.ChatParticipants {
		ws.ChatUsers[chat.ID][participant.ID] += len(ws.UserConnections[participant.ID])
	}
	ws.mu.Unlock()

	ws.SendResponse(cw, wstypes.TypeChatCreated, http.StatusCreated, chat)
}

/*
	──────────────────────────────────────────────────────────────
	Вспомогательные функции
	──────────────────────────────────────────────────────────────
*/

func (ws *WebSocketManager) SendResponse(cw *wstypes.ConnWriter, resptype string, status int, payload any) {
	if err := wstypes.SendResponse(cw, resptype, status, payload); err != nil {
		ws.logger.Error("Send error", "err", err)
	}
}

func (ws *WebSocketManager) SendError(cw *wstypes.ConnWriter, err error) {
	msg, status := helpers.MapError(err)
	if sendErr := wstypes.SendError(cw, status, msg); sendErr != nil {
		ws.logger.Error("Send error", "err", sendErr)
	}
}
