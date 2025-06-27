package storage

import (
	"sync"

	"github.com/yourusername/TouchlineTactics/internal/domain"
)

type MemoryStore struct {
	Rooms map[string]*domain.Room
	Users map[string]*domain.User
	Mutex sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		Rooms: make(map[string]*domain.Room),
		Users: make(map[string]*domain.User),
	}
}

// Room operations
func (s *MemoryStore) GetRoom(id string) (*domain.Room, bool) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	r, ok := s.Rooms[id]
	return r, ok
}

func (s *MemoryStore) SaveRoom(room *domain.Room) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Rooms[room.ID] = room
}

func (s *MemoryStore) DeleteRoom(id string) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	delete(s.Rooms, id)
}

// User operations
func (s *MemoryStore) GetUser(id string) (*domain.User, bool) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	u, ok := s.Users[id]
	return u, ok
}

func (s *MemoryStore) SaveUser(user *domain.User) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Users[user.ID.String()] = user
}

func (s *MemoryStore) DeleteUser(id string) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	delete(s.Users, id)
}

func (s *MemoryStore) ListRooms() []*domain.Room {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	rooms := make([]*domain.Room, 0, len(s.Rooms))
	for _, room := range s.Rooms {
		rooms = append(rooms, room)
	}
	return rooms
}
