package methods

import "github.com/tousart/messenger/internal/models"

type Method interface {
	ServeMessenger(req models.WSRequest)
}

type MessengerMethod func(req models.WSRequest)

func (m MessengerMethod) ServeMessenger(req models.WSRequest) {
	m(req)
}
