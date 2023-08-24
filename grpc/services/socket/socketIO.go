package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Customize as needed
	},
}

func handleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Client connected")

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}

		fmt.Printf("Received message: %s\n", msg)

		err = conn.WriteMessage(messageType, msg)
		if err != nil {
			fmt.Println("Write error:", err)
			break
		}
	}

	fmt.Println("Client disconnected")
}

func main() {
	c := cors.Default()

	http.Handle("/ws", c.Handler(http.HandlerFunc(handleWebSocketConnection)))

	fmt.Println("WebSocket server is up and running on :5000/ws")
	http.ListenAndServe(":5000", nil)
}
