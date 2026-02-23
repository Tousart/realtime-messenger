package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/tousart/messenger/internal/domain"
	"github.com/tousart/messenger/internal/dto"
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
	cookie, err := r.Cookie(domain.CookieSessionID)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userPayload, err := ap.usersService.ValidateSessionID(r.Context(), cookie.Value)
	if err != nil {
		if errors.Is(err, domain.ErrSessionIDNotExists) {
			http.Error(w, "session id not exists", http.StatusUnauthorized)
			return
		}

		log.Printf("authorization error: validate session id: %v\n", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// upgrade connection to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("messengerWebSocketConnectionHandler error: %v\n", err)
		http.Error(w, "websocket connection failed", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	ap.connectUser(userPayload.UserID, conn)
	defer ap.disconnectUser(userPayload.UserID, conn)

	metadata := Metadata{
		UserID: userPayload.UserID,
	}

	for {
		_, wsRequest, err := conn.ReadMessage()
		if err != nil {
			http.Error(w, "websocket connection closed", http.StatusInternalServerError)
			return
		}

		log.Printf("ws request: %s\n", string(wsRequest))

		var req WebSocketRequest
		if err = json.Unmarshal(wsRequest, &req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if method, ok := ap.messengerMethods[req.Method]; ok {
			method(&metadata, &req)
		} else {
			log.Printf("process websocket method error: %v\n", domain.ErrMethodNoTAllowed)
		}
	}
}

/*

	Обработчик регистрации пользователя

*/

func (ap *API) registerHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	response, err := ap.usersService.RegisterUser(r.Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrUserExists) {
			http.Error(w, "user exists", http.StatusBadRequest)
			return
		}
		log.Printf("register error: %v\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	response.RedirectPath = "/"

	http.SetCookie(w, &http.Cookie{
		Name:     domain.CookieSessionID,
		Value:    response.SessionID,
		Path:     "/",
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("register error: %v\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

/*

	Обработчик аутентификации пользователя

*/

func (ap *API) loginHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	response, err := ap.usersService.LoginUser(r.Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		} else if errors.Is(err, domain.ErrIncorrectPassword) {
			http.Error(w, "incorrect password", http.StatusBadRequest)
			return
		}
		log.Printf("login error: %v\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	response.RedirectPath = "/"

	http.SetCookie(w, &http.Cookie{
		Name:     domain.CookieSessionID,
		Value:    response.SessionID,
		Path:     "/",
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("login error: %v\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
