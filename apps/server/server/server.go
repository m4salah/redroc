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
)

type Server struct {
	mux                 chi.Router
	server              *http.Server
	listenAddr          string
	downloadBackendAddr string
	uploadBackendAddr   string
	searchBackendAddr   string
	connTimeout         time.Duration
	skipGcloudAuth      bool
}

type Options struct {
	Host                string
	Port                int
	DownloadBackendAddr string
	UploadBackendAddr   string
	SearchBackendAddr   string
	ConnTimeout         time.Duration
	SkipGcloudAuth      bool
}

func New(opts Options) *Server {
	address := net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))
	mux := chi.NewMux()
	return &Server{
		downloadBackendAddr: opts.DownloadBackendAddr,
		uploadBackendAddr:   opts.UploadBackendAddr,
		searchBackendAddr:   opts.SearchBackendAddr,
		listenAddr:          address,
		mux:                 mux,
		connTimeout:         opts.ConnTimeout,
		skipGcloudAuth:      opts.SkipGcloudAuth,
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
