package server

import (
	"github.com/m4salah/redroc/restful/handlers"
	"github.com/m4salah/redroc/restful/middleware"
)

func (s *Server) setupRoutes() {
	s.mux.Use(middleware.AcceptContentTypeMiddleware)
	handlers.Home(s.mux)
	handlers.Health(s.mux)
	handlers.Download(s.mux, s.downloadBackendAddr, s.log, s.connTimeout, s.skipGcloudAuth)
	handlers.Upload(s.mux, s.uploadBackendAddr, s.log, s.connTimeout, s.skipGcloudAuth)
	handlers.Search(s.mux, s.searchBackendAddr, s.log, s.connTimeout, s.skipGcloudAuth)
}
