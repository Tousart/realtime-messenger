package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisherRepository struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewRabbitMQPublisherRepository(amqpURL, queueName string) (*RabbitMQPublisherRepository, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	queue, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQPublisherRepository{
		conn:    conn,
		channel: ch,
		queue:   queue,
	}, nil
}

func (p *RabbitMQPublisherRepository) PublishMessage(ctx context.Context, messageBytes []byte) error {
	err := p.channel.Publish(
		"",           // exchange
		p.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        messageBytes,
		})

	if err != nil {
		return fmt.Errorf("rabbitmq: PublishMessage error: %s", err.Error())
	}

	return nil
}

func (p *RabbitMQPublisherRepository) Close() {
	p.channel.Close()
	p.conn.Close()
}
