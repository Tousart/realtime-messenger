package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/tousart/messenger/internal/api/helpers"
	"github.com/tousart/messenger/internal/domain"
	"github.com/tousart/messenger/internal/dto"
	"github.com/tousart/messenger/pkg/apirender"
)

/*
	──────────────────────────────────────────────────────────────
	Websocket handler
	──────────────────────────────────────────────────────────────
*/

func (ap *API) messengerWebSocketConnectionHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(domain.CookieSessionID)
	if err != nil {
		apirender.Error(w, http.StatusUnauthorized, domain.ErrUnauthorized.Error())
		return
	}

	sessionPayload, err := ap.usersUC.ValidateSessionID(r.Context(), cookie.Value)
	if err != nil {
		if errors.Is(err, domain.ErrSessionIDNotExists) {
			ap.renderError(w, err)
			return
		}
		ap.renderError(w, err)
		return
	}

	chats, err := ap.messagesUC.UsersChats(r.Context(), sessionPayload.UserID)
	if err != nil {
		ap.renderError(w, err)
		return
	}

	userPayload := &dto.UserPayload{
		ID:    sessionPayload.UserID,
		Name:  sessionPayload.UserName,
		Chats: chats,
	}

	if err = ap.wsUpgrader.UpgradeConnectionForUser(w, r, nil, userPayload); err != nil {
		ap.renderError(w, err)
		return
	}
}

/*
	──────────────────────────────────────────────────────────────
	Авторизация
	──────────────────────────────────────────────────────────────
*/

func (ap *API) registerHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	response, err := ap.usersUC.RegisterUser(r.Context(), req)
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

func (ap *API) loginHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	response, err := ap.usersUC.LoginUser(r.Context(), req)
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

/*
	──────────────────────────────────────────────────────────────
	Вспомогательные функции
	──────────────────────────────────────────────────────────────
*/

func (ap *API) renderError(w http.ResponseWriter, err error) {
	msg, status := helpers.MapError(err)
	if status == http.StatusInternalServerError {
		ap.logger.Error(err.Error())
	} else {
		ap.logger.Info(err.Error())
	}
	apirender.Error(w, status, msg)
}
