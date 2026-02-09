package api

import (
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/tousart/messenger/internal/models"
	"github.com/tousart/messenger/internal/usecase"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == "http://localhost:8080"
	},
}

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
}

type MessengerMethod func(req models.WSRequest)

func NewAPI(msgsHandlerService usecase.MessagesHandlerService) *API {
	return &API{
		UserConnections:    make(map[int][]*websocket.Conn),
		ChatUsers:          make(map[int]map[int]int),
		UserChat:           make(map[int]int),
		mu:                 &sync.RWMutex{},
		msgsHandlerService: msgsHandlerService,
	}
}

func (ap *API) WithHandlers(r *chi.Mux) {
	r.Route("/", func(r chi.Router) {
		// homepage
		r.Get("/", ap.getHomePageHandler)

		// messenger
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
