package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumerRepository struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewRabbitMQConsumerRepository(amqpAddr, queueName string) (*RabbitMQConsumerRepository, error) {
	conn, err := amqp.Dial(amqpAddr)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq: NewRabbitMQConsumerRepository error: %s", err.Error())
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("rabbitmq: NewRabbitMQConsumerRepository error: %s", err.Error())
	}

	queue, err := ch.QueueDeclare(
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
		return nil, fmt.Errorf("rabbitmq: NewRabbitMQConsumerRepository error: %s", err.Error())
	}

	return &RabbitMQConsumerRepository{
		conn:    conn,
		channel: ch,
		queue:   queue,
	}, nil
}

func (c *RabbitMQConsumerRepository) ConsumeMessages(ctx context.Context) (<-chan amqp.Delivery, error) {
	msgs, err := c.channel.Consume(
		c.queue.Name, // имя очереди
		"",           // consumer tag
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq: ConsumeMessages error: %s", err.Error())
	}

	return msgs, nil
}

func (c *RabbitMQConsumerRepository) Close() {
	c.channel.Close()
	c.conn.Close()
}
