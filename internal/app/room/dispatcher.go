package room

import (
	"encoding/json"

	"github.com/yourusername/TouchlineTactics/internal/domain"
)

type ClientConn interface {
	Send([]byte)
	ID() string
}

type IncomingEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventDispatcher struct {
	Handler *RoomEventHandler
}

func NewEventDispatcher(handler *RoomEventHandler) *EventDispatcher {
	return &EventDispatcher{Handler: handler}
}

func (d *EventDispatcher) Dispatch(client ClientConn, message []byte) {
	var event IncomingEvent
	if err := json.Unmarshal(message, &event); err != nil {
		return // Optionally: send error to client
	}

	switch event.Type {
	case string(EventCreateRoom):
		var payload CreateRoomPayload
		if err := json.Unmarshal(event.Payload, &payload); err == nil {
			d.Handler.HandleCreateRoom(client, payload)
		}
	case string(EventJoinRoom):
		var payload JoinRoomPayload
		if err := json.Unmarshal(event.Payload, &payload); err == nil {
			d.Handler.HandleJoinRoom(client, payload)
		}
	case string(EventSetSettings):
		var payload domain.RoomSettings
		if err := json.Unmarshal(event.Payload, &payload); err == nil {
			d.Handler.HandleSetSettings(client, payload)
		}
	case string(EventChatMessage):
		var payload ChatMessagePayload
		if err := json.Unmarshal(event.Payload, &payload); err == nil {
			d.Handler.HandleChatMessage(client, payload)
		}
	case string(EventSetReady):
		var payload SetReadyPayload
		if err := json.Unmarshal(event.Payload, &payload); err == nil {
			d.Handler.HandleSetReady(client, payload)
		}
	case string(EventGetChatHistory):
		var payload GetChatHistoryPayload
		if err := json.Unmarshal(event.Payload, &payload); err == nil {
			d.Handler.HandleGetChatHistory(client, payload)
		}
	case string(EventGetRoomAnalytics):
		var payload GetRoomAnalyticsPayload
		if err := json.Unmarshal(event.Payload, &payload); err == nil {
			d.Handler.HandleGetRoomAnalytics(client, payload)
		}
	case string(EventLeaveRoom):
		var payload LeaveRoomPayload
		if err := json.Unmarshal(event.Payload, &payload); err == nil {
			d.Handler.HandleLeaveRoom(client, payload)
		}
	case string(EventKickUser):
		var payload KickUserPayload
		if err := json.Unmarshal(event.Payload, &payload); err == nil {
			d.Handler.HandleKickUser(client, payload)
		}
		// Add more cases for other events
	}
}
