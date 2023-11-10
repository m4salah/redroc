package server

import (
	"github.com/m4salah/redroc/apps/server/handlers"
	"github.com/m4salah/redroc/apps/server/middleware"
)

func (s *Server) setupRoutes() {
	// Middleware
	s.mux.Use(middleware.CorsMiddleware)
	s.mux.Use(middleware.AcceptContentTypeMiddleware)

	// Handlers
	handlers.Home(s.mux)
	handlers.Health(s.mux)
	handlers.Download(s.mux, s.config.DownloadBackendAddr, s.connTimeout, s.skipGcloudAuth)
	handlers.Upload(s.mux, s.config)
	handlers.Search(s.mux, s.config.SearchBackendAddr, s.connTimeout, s.skipGcloudAuth)
	handlers.SocketIO(s.mux)
}
