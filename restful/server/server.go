package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	log                 *zap.Logger
	mux                 chi.Router
	server              *http.Server
	listenAddr          string
	downloadBackendAddr string
	uploadBackendAddr   string
	searchBackendAddr   string
	connTimeout         time.Duration
}

type Options struct {
	Host                string
	Log                 *zap.Logger
	Port                int
	DownloadBackendAddr string
	UploadBackendAddr   string
	SearchBackendAddr   string
	ConnTimeout         time.Duration
}

func New(opts Options) *Server {
	if opts.Log == nil {
		opts.Log = zap.NewNop()
	}
	address := net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))
	mux := chi.NewMux()
	return &Server{
		downloadBackendAddr: opts.DownloadBackendAddr,
		uploadBackendAddr:   opts.UploadBackendAddr,
		searchBackendAddr:   opts.SearchBackendAddr,
		listenAddr:          address,
		log:                 opts.Log,
		mux:                 mux,
		connTimeout:         opts.ConnTimeout,
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
	s.log.Info("Starting", zap.String("address", s.listenAddr))
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error starting server: %w", err)
	}
	return nil
}

// Stop the Server gracefully within the timeout.
func (s *Server) Stop() error {
	s.log.Info("Stopping")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error stopping server: %w", err)
	}

	return nil
}
