package server

import (
	"log"
	"net/http"

	"github.com/m4salah/redroc/restful/handlers"
)

func acceptContentTypeMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("middleware", r.URL)
		h.ServeHTTP(w, r)
	})
}

func (s *Server) setupRoutes() {
	s.mux.Use(acceptContentTypeMiddleware)
	handlers.Home(s.mux)
	handlers.Health(s.mux)
	handlers.Download(s.mux, s.downloadBackendAddr, s.log, s.connTimeout)
	handlers.Upload(s.mux, s.uploadBackendAddr, s.log, s.connTimeout)
	handlers.Search(s.mux, s.searchBackendAddr, s.log, s.connTimeout)
}
