package room

import (
	"github.com/yourusername/TouchlineTactics/internal/domain"
)

func (h *RoomEventHandler) LeaveRoom(user *domain.User) {
	room, ok := h.Store.GetRoom(user.RoomID)
	if !ok {
		return
	}
	room.Mutex.Lock()
	defer room.Mutex.Unlock()
	delete(room.Users, user.ID.String())
	if len(room.Users) == 0 {
		h.Store.DeleteRoom(room.ID)
		return
	}
	if room.HostID == user.ID.String() {
		// Transfer host to another user
		for _, u := range room.Users {
			u.IsHost = true
			room.HostID = u.ID.String()
			break
		}
	}
	h.Store.SaveRoom(room)
	h.Broadcast(room.ID, EventRoomStateUpdate, room)
}
