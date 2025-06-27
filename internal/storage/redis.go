package storage

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/yourusername/TouchlineTactics/internal/domain"
)

type RedisStore struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewRedisStore(addr, password string, db int) *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisStore{
		Client: client,
		Ctx:    context.Background(),
	}
}

// Room operations
func (s *RedisStore) GetRoom(id string) (*domain.Room, bool) {
	val, err := s.Client.Get(s.Ctx, "room:"+id).Result()
	if err != nil {
		return nil, false
	}
	var room domain.Room
	if err := json.Unmarshal([]byte(val), &room); err != nil {
		return nil, false
	}
	return &room, true
}

func (s *RedisStore) SaveRoom(room *domain.Room) {
	b, _ := json.Marshal(room)
	s.Client.Set(s.Ctx, "room:"+room.ID, b, 0)
}

func (s *RedisStore) DeleteRoom(id string) {
	s.Client.Del(s.Ctx, "room:"+id)
}

// User operations
func (s *RedisStore) GetUser(id string) (*domain.User, bool) {
	val, err := s.Client.Get(s.Ctx, "user:"+id).Result()
	if err != nil {
		return nil, false
	}
	var user domain.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, false
	}
	return &user, true
}

func (s *RedisStore) SaveUser(user *domain.User) {
	b, _ := json.Marshal(user)
	s.Client.Set(s.Ctx, "user:"+user.ID.String(), b, 0)
}

func (s *RedisStore) DeleteUser(id string) {
	s.Client.Del(s.Ctx, "user:"+id)
}

// Pub/Sub for distributed events
func (s *RedisStore) PublishEvent(channel string, data interface{}) {
	b, _ := json.Marshal(data)
	s.Client.Publish(s.Ctx, channel, b)
}

func (s *RedisStore) SubscribeEvents(channel string, handler func([]byte)) {
	pubsub := s.Client.Subscribe(s.Ctx, channel)
	ch := pubsub.Channel()
	go func() {
		for msg := range ch {
			handler([]byte(msg.Payload))
		}
	}()
}

func (s *RedisStore) ListRooms() []*domain.Room {
	var rooms []*domain.Room
	iter := s.Client.Scan(s.Ctx, 0, "room:*", 0).Iterator()
	for iter.Next(s.Ctx) {
		val, err := s.Client.Get(s.Ctx, iter.Val()).Result()
		if err != nil {
			continue
		}
		var room domain.Room
		if err := json.Unmarshal([]byte(val), &room); err == nil {
			rooms = append(rooms, &room)
		}
	}
	return rooms
}
