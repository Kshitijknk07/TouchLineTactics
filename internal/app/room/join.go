package room

import (
	"errors"

	"github.com/yourusername/TouchlineTactics/internal/domain"
)

func (h *RoomEventHandler) JoinRoom(user *domain.User, room *domain.Room, password string) error {
	if room.Settings.Private && room.Settings.Password != password {
		return errors.New("invalid password or private room")
	}
	if IsRoomAtCapacity(room) {
		return errors.New("room is at capacity")
	}
	room.Mutex.Lock()
	defer room.Mutex.Unlock()
	room.Users[user.ID.String()] = user
	h.Store.SaveRoom(room)
	return nil
}
