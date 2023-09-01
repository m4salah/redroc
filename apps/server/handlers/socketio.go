package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/m4salah/redroc/apps/server/pkg/socketio"
)

func SocketIO(mux chi.Router) {
	manager := socketio.NewManager()
	mux.Get("/ws", manager.ServeWS)
}
