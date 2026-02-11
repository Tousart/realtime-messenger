package rabbitmq

import (
	"github.com/tousart/messenger/internal/domain"
)

func (r *RabbitMQMessagesHandlerRepository) MessagesQueue() (domain.MessagesQueue, error) {
	return domain.MessagesQueue(r.messagesQueue), nil
}
