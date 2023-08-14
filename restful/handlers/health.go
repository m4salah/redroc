package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/m4salah/redroc/restful/middleware"
)

func Health(mux chi.Router) {
	mux.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		acceptType := ctx.Value(middleware.GetAcceptTypeKey) == middleware.JSONType
		fmt.Println("acceptType", acceptType)
		fmt.Fprintf(w, "Healthy!")
	})
}
