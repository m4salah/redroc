package socketio

import (
	"fmt"
	"log/slog"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

type Client struct {
	conn    *websocket.Conn
	manager *Manager

	// egress is used to avoid concurrent read/write to the websocket connection
	egress chan []byte
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		conn:    conn,
		manager: manager,
		egress:  make(chan []byte),
	}
}

func (c *Client) readMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	for {
		mt, payload, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("error reading message", "error", err)
			}
			break
		}

		for wsclient := range c.manager.clients {
			slog.Info("sending message to client", "addr", wsclient.conn.RemoteAddr().String())
			wsclient.egress <- payload
		}
		fmt.Println(mt, string(payload))

	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					slog.Error("error writing close message", "error", err)
				}
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				slog.Error("error writing message", "error", err)
			}
			slog.Info("message sent", slog.String("message", string(message)))
		}
	}
}
