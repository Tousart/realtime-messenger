package wstypes

import "encoding/json"

const (
	TypeError          = "error"
	TypeOk             = "ok"
	TypeMessageNew     = "message.new"
	TypeMessageCreated = "message.created"
	TypeChatCreated    = "chat.created"
)

type Response struct {
	Type    string `json:"type"`
	Payload any    `json:"payload,omitempty"`
	Status  int    `json:"status,omitempty"`
}

type Error struct {
	Msg string `json:"error"`
}

func SendResponse(cw *ConnWriter, resptype string, status int, payload any) error {
	return Send(cw, Response{
		Type:    resptype,
		Payload: payload,
		Status:  status,
	})
}

func SendError(cw *ConnWriter, status int, msg string) error {
	return SendResponse(cw, TypeError, status, Error{Msg: msg})
}

func Send(cw *ConnWriter, payload any) error {
	msg, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return cw.Send(msg)
}
