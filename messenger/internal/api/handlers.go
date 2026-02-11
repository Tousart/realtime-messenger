package api

import (
	"encoding/json"
	"net/http"

	"github.com/tousart/messenger/internal/domain"
)

func (ap *API) registerHandler(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.UserName == "" || req.Password == "" {
		http.Error(w, "username and password are required", http.StatusBadRequest)
		return
	}

	if err := ap.usersService.RegisterUser(r.Context(), &req); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (ap *API) loginHandler(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.UserName == "" || req.Password == "" {
		http.Error(w, "username and password are required", http.StatusBadRequest)
		return
	}

	if err := ap.usersService.LoginUser(r.Context(), &req); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
