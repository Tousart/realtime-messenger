package websocket

import (
	"encoding/json"
	"errors"

	"github.com/tousart/messenger/internal/api/helpers"
	"github.com/tousart/messenger/internal/domain"
)

const (
	TypeError          = "error"
	TypeOk             = "ok"
	TypeMessageNew     = "message.new"
	TypeMessageCreated = "message.created"
	TypeChatCreated    = "chat.created"
)

type WSResponse struct {
	Type    string `json:"type"`
	Payload any    `json:"payload,omitempty"`
	Status  int    `json:"status,omitempty"`
	Error   string `json:"error,omitempty"`
}

func (ws *WebSocketManager) SendError(cw *ConnWriter, err error) {
	if errors.Is(err, domain.ErrInternal) {
		ws.logger.Error("ERROR:", "err", err)
	} else {
		ws.logger.Info("error:", "err", err)
	}
	msg, status := helpers.MapError(err)
	ws.Send(cw, WSResponse{
		Type:   TypeError,
		Status: status,
		Error:  msg,
	})
}

func (ws *WebSocketManager) SendResponse(cw *ConnWriter, respType string, status int, payload any) {
	ws.Send(cw, WSResponse{
		Type:    respType,
		Payload: payload,
		Status:  status,
	})
}

func (ws *WebSocketManager) Send(cw *ConnWriter, response any) {
	msg, err := json.Marshal(response)
	if err != nil {
		ws.logger.Error("websocket: Send:", "err", err)
	}

	if err = cw.Send(msg); err != nil {
		ws.logger.Error("websocket: Send:", "err", err)
	}
}
