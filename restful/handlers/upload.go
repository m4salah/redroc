package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func Upload(mux chi.Router, backendAddr string, log *zap.Logger, backendTimeout time.Duration) {
	mux.Post("/upload", func(w http.ResponseWriter, r *http.Request) {

	})
}
