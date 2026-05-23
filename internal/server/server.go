package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"hookwatch/internal/storage"
)

// Config holds server configuration.
type Config struct {
	Port     int
	MaxItems int
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Port:     8080,
		MaxItems: 100,
	}
}

// Server wraps the HTTP server and its dependencies.
type Server struct {
	cfg    Config
	http   *http.Server
	store  *storage.Store
}

// New creates a new Server from the given config.
func New(cfg Config) *Server {
	store := storage.NewStore(cfg.MaxItems)
	router := NewRouter(store)

	httpSrv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		cfg:   cfg,
		http:  httpSrv,
		store: store,
	}
}

// Start begins listening and blocks until the context is cancelled.
func (s *Server) Start(ctx context.Context) error {
	log.Printf("hookwatch listening on %s", s.http.Addr)

	errCh := make(chan error, 1)
	go func() {
		if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.http.Shutdown(shutCtx)
	}
}
