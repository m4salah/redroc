package server

import (
	"github.com/m4salah/redroc/restful/handlers"
)

func (s *Server) setupRoutes() {
	handlers.Health(s.mux)
	handlers.Download(s.mux, s.downloadBackendAddr, s.log, s.connTimeout)
	handlers.Upload(s.mux, s.downloadBackendAddr, s.log, s.connTimeout)
}
