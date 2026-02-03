package api

import (
	"context"
	"fmt"
	"log"
	"router/internal/usecase"

	amqp "github.com/rabbitmq/amqp091-go"
)

type API struct {
	consumerService usecase.MessagesConsumerService
}

func NewAPI(consumerService usecase.MessagesConsumerService) *API {
	return &API{
		consumerService: consumerService,
	}
}

func (ap *API) ConsumeMessages(ctx context.Context) error {
	messagesChannel, err := ap.consumerService.ConsumeMessages(context.TODO())
	if err != nil {
		return fmt.Errorf("api: ConsumeMessages error: %s", err.Error())
	}

	go ap.consumeMessagesFromChannel(ctx, messagesChannel)

	return nil
}

func (ap *API) consumeMessagesFromChannel(ctx context.Context, msgsChan <-chan amqp.Delivery) {
	errChan := make(chan error)

	for {
		select {
		case msg := <-msgsChan:
			nodes, err := ap.consumerService.GetNodes()

			for _, nodeAddress := range nodes {
				err := ap.consumerService.SendMessageToNode(nodeAddress)
			}

		case err := <-errChan:
			log.Printf("api: consumeMessagesFromChannel error: %s", err)
			break
		case <-ctx.Done():
			break
		}
	}
}
