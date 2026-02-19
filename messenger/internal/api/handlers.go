package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/tousart/messenger/internal/dto"
)

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
