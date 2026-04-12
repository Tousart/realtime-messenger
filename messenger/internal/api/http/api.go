package api

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/tousart/messenger/internal/middleware"
)

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

func (ap *API) WithHandlers(r *chi.Mux) {
	r.Use(middleware.CorsMiddleware)
	r.Use(middleware.LoggingMiddleware(ap.logger))

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", ap.registerHandler)
		r.Post("/login", ap.loginHandler)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.Authorization(ap.usersUC))
		r.Get("/messenger", ap.messengerWebSocketConnectionHandler)
	})
}
