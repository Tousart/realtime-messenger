package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Addr   string
	srv    *http.Server
	router *chi.Mux
}

func NewServer(addr string, r *chi.Mux) *Server {
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	return &Server{
		Addr:   addr,
		srv:    srv,
		router: r,
	}
}

func (s *Server) CreateAndRunServer(ctx context.Context) error {
	log.Printf("server run on %s\n", s.Addr)
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("run server error: %s\n", err.Error())
		return err
	}
	log.Printf("server on %s stopped\n", s.Addr)
	return nil
}

func (s *Server) ShutdownServer(ctx context.Context) error {
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		return err
	}
	log.Println("server shutting down graceful")
	return nil
}
