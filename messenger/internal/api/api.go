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

func (ap *API) WithHandlers(r *chi.Mux, isProd bool) {
	// homepage
	r.Get("/", ap.getHomePageHandler)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", ap.registerHandler)
		r.Get("/login", ap.getLoginPageHandler)
		r.Post("/login", ap.loginHandler)
	})

	// method to upgrade connection to websocket
	r.Get("/messenger", ap.messengerWebSocketConnectionHandler)

	// service handlers
	// r.Group(func(r chi.Router) {
	// 	if isProd {
	// 		r.Use(middleware.Authorization(ap.usersService))
	// 	}
	// })
}

func (ap *API) WithMethods() {
	ap.messengerMethods = map[string]MessengerMethod{
		"send": MessengerMethod(ap.SendMessage),
		"join": MessengerMethod(ap.JoinToChat),
		// "leave": methods.MessengerMethod(LeaveChat),
	}
}
