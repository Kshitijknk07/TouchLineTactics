package room

import (
	"time"

	"github.com/yourusername/TouchlineTactics/internal/storage"
)

const RoomExpiryDuration = 30 * time.Minute

func StartRoomExpiryChecker(store *storage.MemoryStore, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			store.Mutex.Lock()
			now := time.Now()
			for id, room := range store.Rooms {
				room.Mutex.RLock()
				lastActivity := room.LastActivity
				room.Mutex.RUnlock()
				if now.Sub(lastActivity) > RoomExpiryDuration {
					delete(store.Rooms, id)
				}
			}
			store.Mutex.Unlock()
		}
	}()
}
