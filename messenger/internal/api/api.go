package api

import (
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/tousart/messenger/internal/usecase"
)

type API struct {
	// methods that are processed by websocket connection
	messengerMethods map[string]MessengerMethod

	// control users connections
	UserConnections map[int][]*websocket.Conn

	// control reflections chat-user
	ChatUsers map[int]map[int]int

	// control reflection user an his current chat
	UserChat map[int]int

	// for isolated access to maps
	mu *sync.RWMutex

	// publisher service to balance messages
	msgsHandlerService usecase.MessagesHandlerService

	// processing users data
	usersService usecase.UsersService
}

func NewAPI(msgsHandlerService usecase.MessagesHandlerService, usersService usecase.UsersService) *API {
	return &API{
		UserConnections:    make(map[int][]*websocket.Conn),
		ChatUsers:          make(map[int]map[int]int),
		UserChat:           make(map[int]int),
		mu:                 &sync.RWMutex{},
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
