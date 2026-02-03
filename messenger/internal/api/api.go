package api

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/tousart/messenger/internal/api/methods"
	"github.com/tousart/messenger/internal/models"
	"github.com/tousart/messenger/internal/usecase"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type API struct {
	// methods that are processed by websocket connection
	messengerMethods map[string]methods.MessengerMethod

	// control users connections
	UserConnections map[int][]*websocket.Conn

	// control reflections chat-user
	ChatUsers map[int]map[int]int

	// control reflection user an his current chat
	UserChat map[int]int

	// for isolated access to maps
	mu *sync.RWMutex

	// publisher service to balance messages
	publisherService usecase.MessagesPublisherService
}

func NewAPI(publisherService usecase.MessagesPublisherService) *API {
	return &API{
		UserConnections:  make(map[int][]*websocket.Conn),
		ChatUsers:        make(map[int]map[int]int),
		UserChat:         make(map[int]int),
		mu:               &sync.RWMutex{},
		publisherService: publisherService,
	}
}

func (ap *API) WithHandlers(r *chi.Mux) {
	r.Route("/", func(r chi.Router) {
		r.Get("/", ap.getHomePageHandler)
		r.Get("/messenger", ap.messengerWebSocketConnectionHandler)
	})
}

func (ap *API) WithMethods() {
	ap.messengerMethods = map[string]methods.MessengerMethod{
		"send": methods.MessengerMethod(ap.SendMessage),
		"join": methods.MessengerMethod(ap.JoinToChat),
		// "leave": methods.MessengerMethod(LeaveChat),
	}
}

func (ap *API) messengerWebSocketConnectionHandler(w http.ResponseWriter, r *http.Request) {
	// TODO...
	// authorization
	// get userID
	userID := UserID

	// upgrade connection to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("messengerWebSocketConnectionHandler error: %s\n", err.Error())
		return
	}
	defer conn.Close()

	// ctx and errors channel for correct shutdown connection
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error)

	// connect and disconnect user
	// add users connection and chats to maps
	ap.connectUser(userID, conn)
	defer ap.disconnectUser(userID, conn)

	// handle requests on websocket connection
	go ap.wsRequestHandler(ctx, conn, userID, errChan)

	select {
	case <-ctx.Done():
	case err := <-errChan:
		log.Printf("messengerWebSocketConnectionHandler error: %s\n", err.Error())
	}
}

func (ap *API) wsRequestHandler(ctx context.Context, conn *websocket.Conn, userID int, errChan chan<- error) {
	for {
		_, wsRequest, err := conn.ReadMessage()

		if err != nil {
			select {
			case <-ctx.Done():
			case errChan <- err:
			}
			return
		}

		log.Printf("ws request: %s\n", string(wsRequest))

		var wsReq models.WebSocketMessageRequest
		if err := json.Unmarshal(wsRequest, &wsReq); err != nil {
			select {
			case <-ctx.Done():
			case errChan <- err:
			}
			return
		}

		if method, ok := ap.messengerMethods[wsReq.Method]; ok {
			req := models.WSRequest{
				UserID: userID,
				Data:   wsReq.Data,
			}
			go method(req)
		} else {
			log.Printf("wsRequestHandler error: %v", errors.New("method not allowed"))
		}
	}
}

func (ap *API) connectUser(userID int, conn *websocket.Conn) {
	// TODO...
	// get users chats
	chats := Chats

	ap.mu.Lock()
	defer ap.mu.Unlock()

	// user - connection
	if ap.UserConnections[userID] == nil {
		ap.UserConnections[userID] = []*websocket.Conn{conn}
	} else {
		ap.UserConnections[userID] = append(ap.UserConnections[userID], conn)
	}

	// chatID - user
	for _, chatID := range chats {
		if _, ok := ap.ChatUsers[chatID]; !ok {
			ap.ChatUsers[chatID] = make(map[int]int)
		}
		ap.ChatUsers[chatID][userID]++
	}
}

func (ap *API) disconnectUser(userID int, conn *websocket.Conn) {
	// TODO...
	// get users chats
	chats := Chats

	ap.mu.Lock()
	defer ap.mu.Unlock()

	// user - connection
	connections := ap.UserConnections[userID]

	if len(connections) == 1 {
		delete(ap.UserConnections, userID)
	} else {
		for i, c := range connections {
			if c == conn {
				ap.UserConnections[userID] = append(connections[:i], connections[i+1:]...)
				break
			}
		}
	}

	// chatID - user
	for _, chatID := range chats {
		if ap.ChatUsers[chatID][userID] == 1 {
			delete(ap.ChatUsers[chatID], userID)

			if len(ap.ChatUsers[chatID]) == 0 {
				delete(ap.ChatUsers, chatID)
			}
		} else {
			ap.ChatUsers[chatID][userID]--
		}
	}
}
