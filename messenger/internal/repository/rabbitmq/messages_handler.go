package rabbitmq

import (
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQMessagesHandlerRepository struct {
	channel       *amqp.Channel
	messagesQueue <-chan amqp.Delivery

	// other nodes queues
	declaredQueues map[string]bool
	mu             *sync.RWMutex
}

func NewRabbitMQMessagesHandlerRepository(channel *amqp.Channel, messagesQueue <-chan amqp.Delivery, queueName string) (*RabbitMQMessagesHandlerRepository, error) {
	_, err := channel.QueueDeclare(
		queueName, // имя очереди
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // аргументы
	)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq: NewRabbitMQMessagesHandlerRepository error: %s", err.Error())
	}
	declaredQueues := map[string]bool{queueName: true}

	return &RabbitMQMessagesHandlerRepository{
		channel:        channel,
		messagesQueue:  messagesQueue,
		declaredQueues: declaredQueues,
		mu:             &sync.RWMutex{},
	}, nil
}
