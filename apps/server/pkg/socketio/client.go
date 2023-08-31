package socketio

import (
	"fmt"
	"log"

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
				log.Println("error reading message", err)
			}
			break
		}

		for wsclient := range c.manager.clients {
			fmt.Println("sending message to client", wsclient.conn.RemoteAddr().String())
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
					log.Println("error writing close message", err)
				}
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("error writing message", err)
			}
			log.Println("message sent", string(message))
		}
	}
}
