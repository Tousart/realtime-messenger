package repository

import (
	"context"
)

type MessagesPublisherRepository interface {
	PublishMessage(ctx context.Context, messageBytes []byte) error
}
