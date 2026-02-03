package usecase

import (
	"context"

	"github.com/tousart/messenger/internal/models"
)

type MessagesPublisherService interface {
	PublishMessage(ctx context.Context, message models.Message) error
}
