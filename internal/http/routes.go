package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/TouchlineTactics/internal/app/room"
	"github.com/yourusername/TouchlineTactics/internal/ws"
)

func SetupRoutes(app *fiber.App, hub *ws.Hub, dispatcher *room.EventDispatcher) {
	app.Get("/ws", ws.WebSocketHandler(hub, dispatcher))
	// Add HTTP endpoints here, e.g.:
	// app.Get("/rooms", ListRoomsHandler)
}
