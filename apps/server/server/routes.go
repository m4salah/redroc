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
	handlers.Download(s.mux, s.downloadBackendAddr, s.connTimeout, s.skipGcloudAuth)
	handlers.Upload(s.mux, s.uploadBackendAddr, s.connTimeout, s.skipGcloudAuth)
	handlers.Search(s.mux, s.searchBackendAddr, s.connTimeout, s.skipGcloudAuth)
	handlers.SocketIO(s.mux)
}
