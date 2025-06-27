package room

import (
	"github.com/google/uuid"
	"github.com/yourusername/TouchlineTactics/internal/domain"
)

func GenerateReconnectToken() string {
	return uuid.NewString()
}

func AssociateReconnectToken(store Store, userID, token string) {
	user, ok := store.GetUser(userID)
	if !ok {
		return
	}
	user.ReconnectToken = token
	store.SaveUser(user)
}

func ValidateReconnectToken(store Store, token string) (*domain.User, bool) {
	// For MemoryStore, scan all users; for RedisStore, you may need a secondary index or scan
	// Here, we assume a scan for simplicity
	// (In production, optimize for Redis with a token->userID map)
	return nil, false // TODO: Implement for RedisStore
}
