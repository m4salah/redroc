package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/m4salah/redroc/apps/server/pkg/socketio"
	"go.uber.org/zap"
)

func SocketIO(mux chi.Router, logger *zap.Logger) {
	manager := socketio.NewManager()
	mux.Get("/ws", manager.ServeWS)
}
