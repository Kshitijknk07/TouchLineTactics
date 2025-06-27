package ws

import (
	"github.com/yourusername/TouchlineTactics/internal/app/room"
	"github.com/yourusername/TouchlineTactics/internal/storage"
)

func OnDisconnect(store *storage.MemoryStore, handler *room.RoomEventHandler, userID string) {
	user, ok := store.GetUser(userID)
	if !ok {
		return
	}
	handler.LeaveRoom(user)
	store.DeleteUser(userID)
}
