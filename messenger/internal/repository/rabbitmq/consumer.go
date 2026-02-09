package rabbitmq

import (
	"github.com/tousart/messenger/internal/models"
)

func (r *RabbitMQMessagesHandlerRepository) MessagesQueue() (models.MessagesQueue, error) {
	return models.MessagesQueue(r.messagesQueue), nil
}
