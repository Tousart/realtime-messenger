package usecase

import (
	"context"

	"github.com/tousart/messenger/internal/dto"
)

type MessagesHandlerService interface {
	PublishMessageToQueues(ctx context.Context, message dto.SendMessageWSRequest) error
	SendMessageToUsersConnections(ctx context.Context, input dto.ConsumingMessage) error
	AddQueueToChat(ctx context.Context, input dto.ChatWSRequest) error
	RemoveQueueFromChat(ctx context.Context, input dto.ChatWSRequest) error
}

// WebSocketManager Interface - для того, чтобы usecase вызывал методы структур, хранящих websocket-сосединения.
// Неважно, что это будут за структуры (usecase о них не знает).
// Обманка в том, что менеджер websocket-соединений находится в слое api, про который usecase знать не должен,
// но он обращается к объекту, реализующему интерфейс websocket-менеджера,
// поэтому неважно из api он или нет - слои друг о друге не знают.
type WebSocketManager interface {
	SendMessageToUsersConnections(ctx context.Context, message dto.ConsumingMessage) error
}
