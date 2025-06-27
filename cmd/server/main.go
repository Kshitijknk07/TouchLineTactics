package main

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/TouchlineTactics/internal/app/auction"
	"github.com/yourusername/TouchlineTactics/internal/app/room"
	apphttp "github.com/yourusername/TouchlineTactics/internal/http"
	"github.com/yourusername/TouchlineTactics/internal/storage"
	"github.com/yourusername/TouchlineTactics/internal/ws"
)

func main() {
	app := fiber.New()

	hub := ws.NewHub()
	go hub.Run()

	var store room.Store
	useRedis := os.Getenv("USE_REDIS") == "1"
	var redisStore *storage.RedisStore
	if useRedis {
		redisStore = storage.NewRedisStore("localhost:6379", "", 0)
		store = redisStore
	} else {
		store = storage.NewMemoryStore()
	}

	roomService := room.NewRoomService()

	// Map of roomID to connected clients (for broadcast)
	var roomClients = make(map[string]map[string]*ws.Client)
	var mu sync.RWMutex

	// Broadcast function
	broadcast := func(roomID string, eventType room.EventType, data interface{}) {
		msg, _ := json.Marshal(map[string]interface{}{
			"type":    eventType,
			"payload": data,
		})
		if useRedis {
			redisStore.PublishEvent("room:"+roomID, msg)
		}
		mu.RLock()
		clients, ok := roomClients[roomID]
		mu.RUnlock()
		if !ok {
			return
		}
		for _, client := range clients {
			client.Send(msg)
		}
	}

	handler := &room.RoomEventHandler{
		Store:       store,
		RoomService: roomService,
		Broadcast:   broadcast,
	}
	dispatcher := room.NewEventDispatcher(handler)

	// Subscribe to Redis pub/sub for distributed events
	if useRedis {
		redisStore.SubscribeEvents("room:*", func(msg []byte) {
			var event map[string]interface{}
			_ = json.Unmarshal(msg, &event)
			roomID, _ := event["roomID"].(string)
			mu.RLock()
			clients, ok := roomClients[roomID]
			mu.RUnlock()
			if ok {
				for _, client := range clients {
					client.Send(msg)
				}
			}
		})
	}

	// Register clients to roomClients map on connect
	apphttp.SetupRoutes(app, hub, dispatcher)

	// Adapt broadcast to match auction.NewAuctionService signature
	auctionBroadcast := func(roomID string, eventType interface{}, data interface{}) {
		broadcast(roomID, eventType.(room.EventType), data)
	}

	auctionService := auction.NewAuctionService(auctionBroadcast, redisStore)
	auctionHandler := &auction.AuctionEventHandler{Auction: auctionService}
	handler.AuctionHandler = auctionHandler

	app.Listen(":8080")
}
