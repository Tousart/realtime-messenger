package api

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/tousart/messenger/internal/middleware"
)

// func corsMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		origin := r.Header.Get("Origin")
// 		if origin == "http://localhost:3000" {
// 			w.Header().Set("Access-Control-Allow-Origin", origin)
// 			w.Header().Set("Access-Control-Allow-Credentials", "true")
// 			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
// 			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
// 		}
// 		if r.Method == http.MethodOptions {
// 			w.WriteHeader(http.StatusNoContent)
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }

type API struct {
	// WebSocketUpgrader
	wsUpgrader WebSocketUpgrader

	// messages usecase
	msgsUC MessagesUsecase

	// processing users data
	usersUC UsersUsecase

	// logger
	logger *slog.Logger
}

func NewAPI(wsUpgrader WebSocketUpgrader, msgsUC MessagesUsecase, usersUC UsersUsecase, logger *slog.Logger) *API {
	return &API{
		wsUpgrader: wsUpgrader,
		msgsUC:     msgsUC,
		usersUC:    usersUC,
		logger:     logger,
	}
}

func (ap *API) WithHandlers(r *chi.Mux, isProd bool) {
	// r.Use(corsMiddleware)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", ap.registerHandler)
		r.Post("/login", ap.loginHandler)
	})

	// method to upgrade connection to websocket — requires authorization
	r.Group(func(r chi.Router) {
		r.Use(middleware.Authorization(ap.usersUC))
		r.Get("/messenger", ap.messengerWebSocketConnectionHandler)
	})
}
