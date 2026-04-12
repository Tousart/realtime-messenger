package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	srv    *http.Server
	logger *slog.Logger
}

func NewServer(host string, port int, r *chi.Mux, logger *slog.Logger) *Server {
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: r,
	}

	return &Server{
		srv:    srv,
		logger: logger,
	}
}

func (s *Server) CreateAndRunServer(ctx context.Context) error {
	s.logger.Info(fmt.Sprintf("server run on %s", s.srv.Addr))

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Error("run server error", "err", err)
		return err
	}

	s.logger.Info(fmt.Sprintf("server on %s stopped", s.srv.Addr))

	return nil
}

func (s *Server) ShutdownServer(ctx context.Context) error {
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	s.logger.Info("server shutting down graceful")

	return nil
}
