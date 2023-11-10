package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/m4salah/redroc/apps/server/types"
)

type Server struct {
	mux            chi.Router
	server         *http.Server
	listenAddr     string
	connTimeout    time.Duration
	skipGcloudAuth bool
	config         types.Config
}

type Options struct {
	ConnTimeout    time.Duration
	SkipGcloudAuth bool
	ServerConfig   types.Config
}

func New(opts Options) *Server {
	address := net.JoinHostPort(opts.ServerConfig.Host, strconv.Itoa(opts.ServerConfig.Port))
	mux := chi.NewMux()
	return &Server{
		listenAddr:     address,
		mux:            mux,
		connTimeout:    opts.ConnTimeout,
		skipGcloudAuth: opts.SkipGcloudAuth,
		config:         opts.ServerConfig,
		server: &http.Server{
			Addr:              address,
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       5 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	s.setupRoutes()
	slog.Info("Starting", slog.String("address", s.listenAddr))
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error starting server: %w", err)
	}
	return nil
}

// Stop the Server gracefully within the timeout.
func (s *Server) Stop() error {
	slog.Info("Stopping")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error stopping server: %w", err)
	}

	return nil
}
