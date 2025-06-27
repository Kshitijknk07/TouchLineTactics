package ws

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	SendChan chan []byte
	IDValue  string
}

func (c *Client) Send(msg []byte) {
	c.SendChan <- msg
}

func (c *Client) ID() string {
	return c.IDValue
}
