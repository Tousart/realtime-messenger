package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/tousart/messenger/internal/dto"
	"github.com/tousart/messenger/internal/middleware"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == "http://localhost:8080"
	},
}

/*

	Обработчик улучшения соединения до WebSocket

*/

func (ap *API) messengerWebSocketConnectionHandler(w http.ResponseWriter, r *http.Request) {
	// build
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Printf("authorization error: get cookie: %v\n", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	userPayload, err := ap.usersService.ValidateSessionID(r.Context(), cookie.Value)
	if err != nil {
		log.Printf("authorization error: validate session id: %v\n", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	// build

	userPayload, ok := r.Context().Value(middleware.ContextKeyAuthMetadata).(*dto.UserPayload)
	if !ok || userPayload == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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

	// connect and disconnect user: add users connection and chats to maps
	ap.connectUser(userPayload.UserID, conn)
	defer ap.disconnectUser(userPayload.UserID, conn)

	// handle requests on websocket connection
	metadata := Metadata{
		UserID: userPayload.UserID,
	}
	go ap.wsRequestHandler(ctx, conn, &metadata, errChan)

	select {
	case <-ctx.Done():
	case err := <-errChan:
		log.Printf("messengerWebSocketConnectionHandler error: %s\n", err.Error())
	}
}

func (ap *API) registerHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ошибка при регистрации: %v\n", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	response, err := ap.usersService.RegisterUser(r.Context(), req)
	if err != nil {
		log.Printf("ошибка при регистрации: %v\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("ошибка при регистрации: %v\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (ap *API) loginHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	response, err := ap.usersService.LoginUser(r.Context(), req)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
