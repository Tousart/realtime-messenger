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

func (s *Server) CreateAndRunServer(ctx context.Context, wg *sync.WaitGroup) {
	errChan := make(chan error, 1)

	wg.Go(func() {
		s.runServer(errChan)
	})

	wg.Go(func() {
		s.shutdownServer(ctx, errChan)
	})
}

func (s *Server) runServer(errChan chan error) {
	log.Printf("server run on %s\n", s.Addr)

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("run server error: %s\n", err.Error())
		errChan <- err
		return
	}

	log.Printf("server run on %s stopped\n", s.Addr)
}

func (s *Server) shutdownServer(ctx context.Context, errChan chan error) {
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
