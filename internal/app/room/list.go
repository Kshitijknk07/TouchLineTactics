package room

type RoomListItem struct {
	ID       string `json:"id"`
	Host     string `json:"host"`
	Status   string `json:"status"`
	NumUsers int    `json:"numUsers"`
	MaxUsers int    `json:"maxUsers"`
	Private  bool   `json:"private"`
}

func (h *RoomEventHandler) ListRooms() []RoomListItem {
	var rooms []RoomListItem
	for _, room := range h.Store.ListRooms() {
		room.Mutex.RLock()
		if !room.Settings.Private && !IsRoomAtCapacity(room) {
			rooms = append(rooms, RoomListItem{
				ID:       room.ID,
				Host:     room.HostID,
				Status:   string(room.Status),
				NumUsers: len(room.Users),
				MaxUsers: room.Settings.MaxUsers,
				Private:  room.Settings.Private,
			})
		}
		room.Mutex.RUnlock()
	}
	return rooms
}
