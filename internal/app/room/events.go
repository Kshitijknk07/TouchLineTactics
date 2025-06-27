package room

import (
	"encoding/json"
	"time"

	"github.com/yourusername/TouchlineTactics/internal/app/auction"
	"github.com/yourusername/TouchlineTactics/internal/domain"
)

type EventType string

const (
	EventCreateRoom       EventType = "createRoom"
	EventJoinRoom         EventType = "joinRoom"
	EventLeaveRoom        EventType = "leaveRoom"
	EventRoomStateUpdate  EventType = "roomStateUpdate"
	EventChatMessage      EventType = "chatMessage"
	EventSetReady         EventType = "setReady"
	EventUserAction       EventType = "userAction"
	EventKickUser         EventType = "kickUser"
	EventTransferHost     EventType = "transferHost"
	EventStartPhase       EventType = "startPhase"
	EventSetSettings      EventType = "setSettings"
	EventListRooms        EventType = "listRooms"
	EventGetChatHistory   EventType = "getChatHistory"
	EventGetRoomAnalytics EventType = "getRoomAnalytics"
)

type CreateRoomPayload struct {
	Username string                 `json:"username"`
	RoomID   string                 `json:"roomId"`
	Password string                 `json:"password,omitempty"`
	Private  bool                   `json:"private,omitempty"`
	Settings map[string]interface{} `json:"settings,omitempty"`
}

type JoinRoomPayload struct {
	RoomID         string `json:"roomId"`
	Username       string `json:"username"`
	Password       string `json:"password,omitempty"`
	ReconnectToken string `json:"reconnectToken,omitempty"`
}

type TransferHostPayload struct {
	NewHostID string `json:"newHostId"`
}

type StartPhasePayload struct {
	Phase string `json:"phase"`
}

type ChatMessagePayload struct {
	Message string `json:"message"`
}

type SetReadyPayload struct {
	Ready bool `json:"ready"`
}

type GetChatHistoryPayload struct{}

type GetRoomAnalyticsPayload struct{}

type RoomAnalytics struct {
	RoomID        string `json:"roomId"`
	TotalMessages int    `json:"totalMessages"`
	TotalUsers    int    `json:"totalUsers"`
}

type LeaveRoomPayload struct{}

type Store interface {
	GetRoom(string) (*domain.Room, bool)
	SaveRoom(*domain.Room)
	DeleteRoom(string)
	GetUser(string) (*domain.User, bool)
	SaveUser(*domain.User)
	DeleteUser(string)
	ListRooms() []*domain.Room
}

type RoomEventHandler struct {
	Store          Store
	RoomService    *RoomService
	Broadcast      func(roomID string, eventType EventType, data interface{})
	AuctionHandler *auction.AuctionEventHandler
}

// Method signatures for event handling
func (h *RoomEventHandler) HandleCreateRoom(client ClientConn, payload CreateRoomPayload) {}
func (h *RoomEventHandler) HandleJoinRoom(client ClientConn, payload JoinRoomPayload)     {}

func (h *RoomEventHandler) HandleTransferHost(client ClientConn, payload TransferHostPayload) {
	user, ok := h.Store.GetUser(client.ID())
	if !ok {
		return
	}
	room, ok := h.Store.GetRoom(user.RoomID)
	if !ok {
		return
	}
	if room.HostID != user.ID.String() {
		return // Only host can transfer
	}
	room.HostID = payload.NewHostID
	for _, u := range room.Users {
		u.IsHost = (u.ID.String() == payload.NewHostID)
	}
	h.Store.SaveRoom(room)
	h.Broadcast(room.ID, EventRoomStateUpdate, room)
}

func (h *RoomEventHandler) HandleStartPhase(client ClientConn, payload StartPhasePayload) {
	user, ok := h.Store.GetUser(client.ID())
	if !ok {
		return
	}
	room, ok := h.Store.GetRoom(user.RoomID)
	if !ok {
		return
	}
	if room.HostID != user.ID.String() {
		return // Only host can start phase
	}
	room.Status = domain.RoomStatus(payload.Phase)
	h.Store.SaveRoom(room)
	h.Broadcast(room.ID, EventRoomStateUpdate, room)
}

func (h *RoomEventHandler) HandleSetSettings(client ClientConn, payload domain.RoomSettings) {
	user, ok := h.Store.GetUser(client.ID())
	if !ok {
		return
	}
	room, ok := h.Store.GetRoom(user.RoomID)
	if !ok {
		return
	}
	if room.HostID != user.ID.String() {
		return // Only host can update settings
	}
	room.Settings = payload
	h.Store.SaveRoom(room)
	h.Broadcast(room.ID, EventRoomStateUpdate, room)
}

func (h *RoomEventHandler) HandleChatMessage(client ClientConn, payload ChatMessagePayload) {
	user, ok := h.Store.GetUser(client.ID())
	if !ok {
		return
	}
	room, ok := h.Store.GetRoom(user.RoomID)
	if !ok {
		return
	}
	msg := domain.ChatMessage{
		UserID:    user.ID.String(),
		Username:  user.Username,
		Message:   payload.Message,
		Timestamp: time.Now(),
	}
	room.Mutex.Lock()
	room.Chat = append(room.Chat, msg)
	room.Mutex.Unlock()
	h.Store.SaveRoom(room)
	h.Broadcast(room.ID, EventChatMessage, msg)
}

func (h *RoomEventHandler) HandleSetReady(client ClientConn, payload SetReadyPayload) {
	user, ok := h.Store.GetUser(client.ID())
	if !ok {
		return
	}
	user.Ready = payload.Ready
	h.Store.SaveUser(user)
	room, ok := h.Store.GetRoom(user.RoomID)
	if !ok {
		return
	}
	h.Broadcast(room.ID, EventRoomStateUpdate, room)
}

func (h *RoomEventHandler) HandleGetChatHistory(client ClientConn, payload GetChatHistoryPayload) {
	user, ok := h.Store.GetUser(client.ID())
	if !ok {
		return
	}
	room, ok := h.Store.GetRoom(user.RoomID)
	if !ok {
		return
	}
	client.Send(mustMarshal(map[string]interface{}{
		"type":    EventGetChatHistory,
		"payload": room.Chat,
	}))
}

func (h *RoomEventHandler) HandleGetRoomAnalytics(client ClientConn, payload GetRoomAnalyticsPayload) {
	user, ok := h.Store.GetUser(client.ID())
	if !ok {
		return
	}
	room, ok := h.Store.GetRoom(user.RoomID)
	if !ok {
		return
	}
	analytics := RoomAnalytics{
		RoomID:        room.ID,
		TotalMessages: len(room.Chat),
		TotalUsers:    len(room.Users),
	}
	client.Send(mustMarshal(map[string]interface{}{
		"type":    EventGetRoomAnalytics,
		"payload": analytics,
	}))
}

func (h *RoomEventHandler) HandleLeaveRoom(client ClientConn, payload LeaveRoomPayload) {
	user, ok := h.Store.GetUser(client.ID())
	if !ok {
		return
	}
	h.LeaveRoom(user)
	h.Store.DeleteUser(user.ID.String())
}

func (h *RoomEventHandler) HandleKickUser(client ClientConn, payload KickUserPayload) {
	host, ok := h.Store.GetUser(client.ID())
	if !ok {
		return
	}
	h.KickUser(host, payload)
}

func mustMarshal(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

// ... Add other event handler methods ...
