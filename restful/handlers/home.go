package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Home(mux chi.Router) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<h1>Welcome to Redroc</h1>")
	})
}
