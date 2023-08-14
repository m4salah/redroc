package middleware

import (
	"context"
	"log"
	"net/http"
)

type ContextAcceptTypeKey string

const JSON = "application/json"

const GetAcceptTypeKey ContextAcceptTypeKey = "GetWithType"

const (
	JSONType = iota
	HTMLType
)

func AcceptContentTypeMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept")
		log.Println("middleware", r.URL)
		log.Println("accept", accept)
		switch accept {
		case JSON:
			ctx := context.WithValue(r.Context(), GetAcceptTypeKey, JSONType)
			r = r.WithContext(ctx)
		default:
			ctx := context.WithValue(r.Context(), GetAcceptTypeKey, HTMLType)
			r = r.WithContext(ctx)
		}
		h.ServeHTTP(w, r)
	})
}
