package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Addr   string
	srv    *http.Server
	router *chi.Mux
	Wg     *sync.WaitGroup
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
		Wg:     &sync.WaitGroup{},
	}
}

func (s *Server) CreateAndRunServer(ctx context.Context) {
	errChan := make(chan error, 1)

	s.Wg.Add(1)
	go s.runServer(errChan)

	s.Wg.Add(1)
	go s.shutdownServer(ctx, errChan)
}

func (s *Server) runServer(errChan chan error) {
	defer s.Wg.Done()

	log.Printf("server run on %s\n", s.Addr)

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("run server error: %s\n", err.Error())
		errChan <- err
		return
	}

	log.Printf("server run on %s stopped\n", s.Addr)
}

func (s *Server) shutdownServer(ctx context.Context, errChan chan error) {
	defer s.Wg.Done()

	select {
	case <-ctx.Done():
	case <-errChan:
		return
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		fmt.Fprintf(os.Stderr, "error shutting down server: %s", err.Error())
	}

	log.Println("server shutting down graceful")
}
