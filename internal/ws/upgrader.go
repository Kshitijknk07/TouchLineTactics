package ws

import (
	"encoding/json"
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

// WebSocketHandlerWithRoomTracking tracks which room a user is in and updates room membership for real-time broadcasting.
func WebSocketHandlerWithRoomTracking(hub *Hub, dispatcher *room.EventDispatcher, userID string, addClientToRoom func(roomID, userID string, client *Client), removeClientFromAllRooms func(userID string)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fasthttpadaptor.NewFastHTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				return
			}
			client := &Client{
				Conn:     conn,
				SendChan: make(chan []byte, 256),
				IDValue:  userID,
			}
			hub.Register <- client
			go client.WritePump()

			var currentRoomID string
			for {
				_, message, err := client.Conn.ReadMessage()
				if err != nil {
					hub.Unregister <- client
					removeClientFromAllRooms(userID)
					client.Conn.Close()
					break
				}
				// Intercept joinRoom/createRoom/leaveRoom
				type incoming struct {
					Type    string          `json:"type"`
					Payload json.RawMessage `json:"payload"`
				}
				var inc incoming
				_ = json.Unmarshal(message, &inc)
				if inc.Type == "joinRoom" || inc.Type == "createRoom" {
					var payload struct {
						RoomID string `json:"roomId"`
					}
					_ = json.Unmarshal(inc.Payload, &payload)
					if payload.RoomID != "" && payload.RoomID != currentRoomID {
						removeClientFromAllRooms(userID)
						addClientToRoom(payload.RoomID, userID, client)
						currentRoomID = payload.RoomID
					}
				}
				if inc.Type == "leaveRoom" {
					removeClientFromAllRooms(userID)
					currentRoomID = ""
				}
				dispatcher.Dispatch(client, message)
			}
		})(c.Context())
		return nil
	}
}
