package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SocketIO(mux chi.Router, logger *zap.Logger) {
	mux.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error("upgrade:", zap.Error(err))
			return
		}
		defer c.Close()
		logger.Info("client connected")
		c.WriteMessage(websocket.TextMessage, []byte("Hello, client!"))
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			logger.Info("recv: ", zap.String("message", string(message)))
		}
	})
}
