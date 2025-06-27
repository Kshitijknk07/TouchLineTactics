package ws

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/websocket"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"github.com/yourusername/TouchlineTactics/internal/app/room"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

func WebSocketHandler(hub *Hub, dispatcher *room.EventDispatcher) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fasthttpadaptor.NewFastHTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				return
			}
			client := &Client{
				Conn:     conn,
				SendChan: make(chan []byte, 256),
				IDValue:  c.Query("userId"),
			}
			hub.Register <- client
			go client.WritePump()
			client.ReadPumpWithDispatcher(hub, dispatcher)
		})(c.Context())
		return nil
	}
}
