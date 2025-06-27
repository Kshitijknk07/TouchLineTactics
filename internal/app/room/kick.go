package room

import (
	"github.com/yourusername/TouchlineTactics/internal/domain"
)

type KickUserPayload struct {
	UserID string `json:"userId"`
}

func (h *RoomEventHandler) KickUser(host *domain.User, payload KickUserPayload) {
	room, ok := h.Store.GetRoom(host.RoomID)
	if !ok {
		return
	}
	if room.HostID != host.ID.String() {
		return // Only host can kick
	}
	room.Mutex.Lock()
	defer room.Mutex.Unlock()
	delete(room.Users, payload.UserID)
	h.Store.SaveRoom(room)
	h.Broadcast(room.ID, EventRoomStateUpdate, room)
}
