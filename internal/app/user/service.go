package user

import (
	"github.com/google/uuid"
	"github.com/yourusername/TouchlineTactics/internal/domain"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) NewUser(username, roomID string, isHost bool) *domain.User {
	return &domain.User{
		ID:       uuid.New(),
		Username: username,
		RoomID:   roomID,
		IsHost:   isHost,
		Ready:    false,
	}
}

func (s *UserService) SetReady(user *domain.User, ready bool) {
	user.Ready = ready
}
