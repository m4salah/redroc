package server

import (
	"github.com/m4salah/redroc/restful/handlers"
	"github.com/m4salah/redroc/restful/middleware"
)

func (s *Server) setupRoutes() {
	// Middleware
	s.mux.Use(middleware.CorsMiddleware)
	s.mux.Use(middleware.AcceptContentTypeMiddleware)

	// Handlers
	handlers.Home(s.mux)
	handlers.Health(s.mux)
	handlers.Download(s.mux, s.downloadBackendAddr, s.log, s.connTimeout, s.skipGcloudAuth)
	handlers.Upload(s.mux, s.uploadBackendAddr, s.log, s.connTimeout, s.skipGcloudAuth)
	handlers.Search(s.mux, s.searchBackendAddr, s.log, s.connTimeout, s.skipGcloudAuth)
	handlers.SocketIO(s.mux, s.log)
}
