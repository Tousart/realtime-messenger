package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/tousart/messenger/internal/usecase"
)

type API struct {
	// methods that are processed by websocket connection
	messengerMethods map[string]MessengerMethod

	// WebSocketManager
	wsManager *WebSocketManager

	// publisher service to balance messages
	msgsHandlerService usecase.MessagesHandlerService

	// processing users data
	usersService usecase.UsersService
}

func NewAPI(wsManager *WebSocketManager, msgsHandlerService usecase.MessagesHandlerService, usersService usecase.UsersService) *API {
	return &API{
		wsManager:          wsManager,
		msgsHandlerService: msgsHandlerService,
		usersService:       usersService,
	}
}

func (ap *API) WithHandlers(r *chi.Mux) {
	r.Route("/", func(r chi.Router) {
		// homepage
		r.Get("/", ap.getHomePageHandler)

		// messenger
		// method to upgrade connection to websocket
		r.Get("/messenger", ap.messengerWebSocketConnectionHandler)

		// authorization
		r.Route("/auth", func(r chi.Router) {
			// r.Get("/", ap.getAuthPage)
			r.Post("/register", ap.registerHandler)
			r.Post("/login", ap.loginHandler)
		})
	})
}

func (ap *API) WithMethods() {
	ap.messengerMethods = map[string]MessengerMethod{
		"send": MessengerMethod(ap.SendMessage),
		"join": MessengerMethod(ap.JoinToChat),
		// "leave": methods.MessengerMethod(LeaveChat),
	}
}
