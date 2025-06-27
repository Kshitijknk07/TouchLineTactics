package ws

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/yourusername/TouchlineTactics/internal/app/room"
)

func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}
		hub.Broadcast <- message
	}
}

func (c *Client) ReadPumpWithDispatcher(hub *Hub, dispatcher *room.EventDispatcher) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}
		dispatcher.Dispatch(c, message)
	}
}

func (c *Client) WritePump() {
	for msg := range c.SendChan {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("write error:", err)
			break
		}
	}
	c.Conn.Close()
}
