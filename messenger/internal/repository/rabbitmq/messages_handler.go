package rabbitmq

import (
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQMessagesHandlerRepository struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	//queue   amqp.Queue
	messagesQueue <-chan amqp.Delivery

	// other nodes queues
	declaredQueues map[string]bool

	mu *sync.RWMutex
}

func NewRabbitMQMessagesHandlerRepository(amqpURL, queueName string) (*RabbitMQMessagesHandlerRepository, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queueName, // имя очереди
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // аргументы
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("rabbitmq: NewRabbitMQMessagesHandlerRepository error: %s", err.Error())
	}
	declaredQueues := map[string]bool{queueName: true}

	msgsQueue, err := ch.Consume(
		queueName, // имя очереди
		"",        // consumer tag
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("rabbitmq: NewRabbitMQMessagesHandlerRepository error: %s", err.Error())
	}

	return &RabbitMQMessagesHandlerRepository{
		conn:           conn,
		channel:        ch,
		messagesQueue:  msgsQueue,
		declaredQueues: declaredQueues,
		mu:             &sync.RWMutex{},
	}, nil
}

func (p *RabbitMQMessagesHandlerRepository) Close() {
	p.channel.Close()
	p.conn.Close()
}
