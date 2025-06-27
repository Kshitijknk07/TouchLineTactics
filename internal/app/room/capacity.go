package room

import "github.com/yourusername/TouchlineTactics/internal/domain"

func IsRoomAtCapacity(room *domain.Room) bool {
	max := room.Settings.MaxUsers
	if max == 0 {
		max = 4 // default max users
	}
	return len(room.Users) >= max
}
