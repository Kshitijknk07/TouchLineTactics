package room

import (
	"errors"

	"github.com/yourusername/TouchlineTactics/internal/domain"
)

type RoomService struct{}

func NewRoomService() *RoomService {
	return &RoomService{}
}

func (s *RoomService) NewRoom(id, hostID string, settings domain.RoomSettings) *domain.Room {
	return &domain.Room{
		ID:       id,
		HostID:   hostID,
		Users:    make(map[string]*domain.User),
		Settings: settings,
		Status:   domain.RoomWaiting,
		Chat:     []domain.ChatMessage{},
	}
}

func (s *RoomService) AddUser(room *domain.Room, user *domain.User) error {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()
	if _, exists := room.Users[user.ID.String()]; exists {
		return errors.New("user already in room")
	}
	room.Users[user.ID.String()] = user
	return nil
}

func (s *RoomService) RemoveUser(room *domain.Room, userID string) {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()
	delete(room.Users, userID)
}

func (s *RoomService) SetHost(room *domain.Room, userID string) {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()
	room.HostID = userID
	if user, ok := room.Users[userID]; ok {
		user.IsHost = true
	}
}

func (s *RoomService) UpdateSettings(room *domain.Room, settings domain.RoomSettings) {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()
	room.Settings = settings
}

func (s *RoomService) AddChatMessage(room *domain.Room, msg domain.ChatMessage) {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()
	room.Chat = append(room.Chat, msg)
}
