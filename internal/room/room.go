package room

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

// User represents a participant in a room.
type User struct {
	Username       string   `json:"username"`
	SocketID       string   `json:"socketId"`
	Team           []string `json:"team"`
	Funds          int      `json:"funds"`
	Ready          bool     `json:"ready"`
	ReconnectToken string   `json:"reconnectToken,omitempty"`
}

// Room statuses
const (
	StatusWaiting    = "WAITING"
	StatusInProgress = "IN_PROGRESS"
	StatusFinished   = "FINISHED"
)

// Room holds the state and users of a room.
type Room struct {
	Status   string                 `json:"status"`
	Users    []User                 `json:"users"`
	Host     string                 `json:"host"`
	Password string                 `json:"password,omitempty"`
	Private  bool                   `json:"private"`
	Settings map[string]interface{} `json:"settings,omitempty"`
}

var ctx = context.Background()

// MaxUsers is the maximum allowed users in a room.
const MaxUsers = 4

// Save stores the room in Redis.
func Save(rdb *redis.Client, roomID string, room Room) error {
	data, err := json.Marshal(room)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, "room:"+roomID, data, 0).Err()
}

// Get retrieves a room from Redis.
func Get(rdb *redis.Client, roomID string) (Room, error) {
	var room Room
	val, err := rdb.Get(ctx, "room:"+roomID).Result()
	if err != nil {
		return room, err
	}
	err = json.Unmarshal([]byte(val), &room)
	return room, err
}

// RemoveUser removes a user by socketID from the room.
func (r *Room) RemoveUser(socketID string) {
	for i, u := range r.Users {
		if u.SocketID == socketID {
			r.Users = append(r.Users[:i], r.Users[i+1:]...)
			break
		}
	}
}

// KickUser removes a user by username from the room.
func (r *Room) KickUser(username string) {
	for i, u := range r.Users {
		if u.Username == username {
			r.Users = append(r.Users[:i], r.Users[i+1:]...)
			break
		}
	}
}

// SetUserReady sets the ready status for a user by socketID.
func (r *Room) SetUserReady(socketID string, ready bool) {
	for i, u := range r.Users {
		if u.SocketID == socketID {
			r.Users[i].Ready = ready
			break
		}
	}
}

// IsHost checks if the given socketID is the host.
func (r *Room) IsHost(socketID string) bool {
	return r.Host == socketID
}

// IsFull returns true if the room is at capacity.
func (r *Room) IsFull() bool {
	return len(r.Users) >= MaxUsers
}

// ListRooms lists all room IDs in Redis.
func ListRooms(rdb *redis.Client) ([]string, error) {
	keys, err := rdb.Keys(ctx, "room:*").Result()
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(keys))
	for _, k := range keys {
		ids = append(ids, k[5:]) // strip "room:" prefix
	}
	return ids, nil
}

// TransferHost sets the host to the given username's socketID.
func (r *Room) TransferHost(username string) {
	for _, u := range r.Users {
		if u.Username == username {
			r.Host = u.SocketID
			break
		}
	}
}

// SetStatus sets the room status.
func (r *Room) SetStatus(status string) {
	r.Status = status
}

// SetSettings sets custom settings for the room.
func (r *Room) SetSettings(settings map[string]interface{}) {
	r.Settings = settings
}

// AppendChatMessage appends a chat message to the room's chat history in Redis.
func AppendChatMessage(rdb *redis.Client, roomID string, msg map[string]interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return rdb.RPush(ctx, "room:"+roomID+":chat", data).Err()
}

// GetChatHistory retrieves the chat history for a room.
func GetChatHistory(rdb *redis.Client, roomID string, limit int64) ([]map[string]interface{}, error) {
	vals, err := rdb.LRange(ctx, "room:"+roomID+":chat", 0, limit-1).Result()
	if err != nil {
		return nil, err
	}
	history := make([]map[string]interface{}, 0, len(vals))
	for _, v := range vals {
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(v), &m); err == nil {
			history = append(history, m)
		}
	}
	return history, nil
}

// IncrementRoomUsage increments a counter for room usage analytics.
func IncrementRoomUsage(rdb *redis.Client, roomID string) {
	rdb.Incr(ctx, "room:"+roomID+":usage")
}

// IncrementUserActivity increments a counter for user activity analytics.
func IncrementUserActivity(rdb *redis.Client, username string) {
	rdb.Incr(ctx, "user:"+username+":activity")
}

// PublishRoomUpdate publishes a room update to a Redis channel.
func PublishRoomUpdate(rdb *redis.Client, roomID string, update map[string]interface{}) error {
	data, err := json.Marshal(update)
	if err != nil {
		return err
	}
	return rdb.Publish(ctx, "room_update:"+roomID, data).Err()
}

// SubscribeRoomUpdates subscribes to room updates for a room.
func SubscribeRoomUpdates(rdb *redis.Client, roomID string) *redis.PubSub {
	return rdb.Subscribe(ctx, "room_update:"+roomID)
}

// FindUserByReconnectToken finds a user in the room by reconnect token.
func (r *Room) FindUserByReconnectToken(token string) *User {
	for i, u := range r.Users {
		if u.ReconnectToken == token {
			return &r.Users[i]
		}
	}
	return nil
}
