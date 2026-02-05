package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"router/internal/models"
	"router/internal/usecase"

	amqp "github.com/rabbitmq/amqp091-go"
)

type API struct {
	consumerService usecase.MessagesConsumerService
	senderService   usecase.NodesSenderService
}

func NewAPI(consumerService usecase.MessagesConsumerService, senderService usecase.NodesSenderService) *API {
	return &API{
		consumerService: consumerService,
		senderService:   senderService,
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
	defer close(errChan)

	for {
		select {
		case msg := <-msgsChan:
			var message models.Message
			msgBytes := msg.Body
			if err := json.Unmarshal(msgBytes, &message); err != nil {
				log.Printf("api: consumeMessagesFromChannel error: %s\n", err.Error())
				continue
			}

			nodes, err := ap.senderService.GetChatNodes(context.Background(), message.ChatID)
			if err != nil {
				log.Printf("api: consumeMessagesFromChannel error: %s\n", err.Error())
				continue
			}

			for _, nodeAddr := range nodes {
				err := ap.senderService.SendMessageToNode(context.Background(), nodeAddr, &message)
				if err != nil {
					log.Printf("api: consumeMessagesFromChannel error: %s\n", err.Error())
				}
				// log.Printf("node address: %s, and message text: %s\n", nodeAddress, message.Text)
			}

		case err := <-errChan:
			log.Printf("api: consumeMessagesFromChannel error: %s", err)
			return
		case <-ctx.Done():
			return
		}
	}
}
