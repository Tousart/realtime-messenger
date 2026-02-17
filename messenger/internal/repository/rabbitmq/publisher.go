package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tousart/messenger/internal/domain"
)

func (r *RabbitMQMessagesHandlerRepository) isDeclared(queueName string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.declaredQueues[queueName]
}

func (r *RabbitMQMessagesHandlerRepository) declareQueue(queueName string) error {
	_, err := r.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("rabbitmq: declareQueue: %s", err.Error())
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.declaredQueues[queueName] = true

	return nil
}

func (r *RabbitMQMessagesHandlerRepository) publishMessage(ctx context.Context, queueName string, messageBytes []byte) error {
	err := r.channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        messageBytes,
		})
	if err != nil {
		return fmt.Errorf("rabbitmq: publishMessage error: %s", err.Error())
	}
	return nil
}

func (r *RabbitMQMessagesHandlerRepository) PublishMessageToQueues(ctx context.Context, queues []string, message *domain.Message) error {
	messageBytes, err := json.Marshal(*message)
	if err != nil {
		return fmt.Errorf("rabbitmq: PublishMessageToQueues error: %s", err.Error())
	}
	for _, queue := range queues {
		if !r.isDeclared(queue) {
			if err := r.declareQueue(queue); err != nil {
				return fmt.Errorf("rabbitmq: PublishMessageToQueues error: %s", err.Error())
			}
		}
		if err := r.publishMessage(ctx, queue, messageBytes); err != nil {
			return fmt.Errorf("rabbitmq: PublishMessageToQueues error: %s", err.Error())
		}
	}
	return nil
}
