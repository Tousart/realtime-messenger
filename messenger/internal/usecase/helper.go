package usecase

import (
	"time"

	"github.com/google/uuid"
)

func generateUUID() uuid.UUID {
	return uuid.New()
}

func timeNowUTC() time.Time {
	return time.Now().UTC().Truncate(time.Second)
}
