package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConnection struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queue     <-chan amqp.Delivery
	queueName string
}

func NewRabbitMQConnection(amqpURL, queueName string) (*RabbitMQConnection, error) {
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
		return nil, fmt.Errorf("rabbitmq: RabbitMQConnection error: %w", err)
	}

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
		return nil, fmt.Errorf("rabbitmq: RabbitMQConnection error: %w", err)
	}

	return &RabbitMQConnection{
		conn:      conn,
		channel:   ch,
		queue:     msgsQueue,
		queueName: queueName,
	}, nil
}

func (rc *RabbitMQConnection) Channel() *amqp.Channel {
	return rc.channel
}

func (rc *RabbitMQConnection) Queue() <-chan amqp.Delivery {
	return rc.queue
}

func (rc *RabbitMQConnection) QueueName() string {
	return rc.queueName
}

func (rc *RabbitMQConnection) Close() {
	rc.channel.Close()
	rc.conn.Close()
}
