package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Health(mux chi.Router) {
	mux.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Healthy!")
	})
}
