package api

import (
	"encoding/json"
	"errors"
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

	sessionPayloadBytes, err := ap.usersUC.ValidateSessionID(r.Context(), cookie.Value)
	if err != nil {
		if errors.Is(err, domain.ErrSessionIDNotExists) {
			ap.renderError(w, err)
			return
		}
		ap.renderError(w, err)
		return
	}
	var sessionPayload dto.SessionPayload
	if err = json.Unmarshal(sessionPayloadBytes, &sessionPayload); err != nil {
		ap.renderError(w, err)
		return
	}

	chats, err := ap.msgsUC.UsersChats(r.Context(), sessionPayload.UserID)
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
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apirender.Error(w, http.StatusBadRequest, "invalid request")
		return
	}
	defer r.Body.Close()

	user, err := ap.usersUC.Register(r.Context(), &req)
	if err != nil {
		ap.renderError(w, err)
		return
	}

	apirender.JSON(w, http.StatusCreated, user)
}

func (ap *API) loginHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apirender.Error(w, http.StatusBadRequest, "invalid request")
		return
	}
	defer r.Body.Close()

	sessionID, err := ap.usersUC.Login(r.Context(), &req)
	if err != nil {
		ap.renderError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     domain.CookieSessionID,
		Value:    sessionID.SessionID,
		Path:     "/",
		HttpOnly: true,
	})

	apirender.JSON(w, http.StatusOK, sessionID)
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
