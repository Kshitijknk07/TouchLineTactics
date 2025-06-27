package domain

import (
	"sync"
	"time"
)

type RoomStatus string

const (
	RoomWaiting    RoomStatus = "WAITING"
	RoomInProgress RoomStatus = "IN_PROGRESS"
	RoomFinished   RoomStatus = "FINISHED"
	RoomPaused     RoomStatus = "PAUSED"
	RoomCancelled  RoomStatus = "CANCELLED"
)

type RoomSettings struct {
	Password string
	Private  bool
	GameMode string
	Timer    int
	MaxUsers int
	Custom   map[string]interface{}
}

type ChatMessage struct {
	UserID    string
	Username  string
	Message   string
	Timestamp time.Time
}

type Room struct {
	ID           string
	HostID       string
	Users        map[string]*User
	Settings     RoomSettings
	Status       RoomStatus
	Chat         []ChatMessage
	LastActivity time.Time
	Mutex        sync.RWMutex
}
