package auction

import (
	"sync"
	"time"

	"github.com/yourusername/TouchlineTactics/internal/domain"
	"github.com/yourusername/TouchlineTactics/internal/storage"
)

type PositionAuction struct {
	Position string
	Players  []domain.Player
	Index    int
}

type AuctionState struct {
	Positions     []PositionAuction
	CurrentPos    int
	CurrentBid    int
	CurrentBidder string
	Timer         *time.Timer
	Mutex         sync.Mutex
}

type AuctionService struct {
	State      map[string]*AuctionState // roomID -> state
	StateMutex sync.Mutex
	Broadcast  func(roomID string, eventType interface{}, data interface{})
	Redis      *storage.RedisStore
}

func (a *AuctionService) Start(d string, players []domain.Player) any {
	panic("unimplemented")
}

func NewAuctionService(broadcast func(roomID string, eventType interface{}, data interface{}), redis *storage.RedisStore) *AuctionService {
	return &AuctionService{
		State:     make(map[string]*AuctionState),
		Broadcast: broadcast,
		Redis:     redis,
	}
}

func (a *AuctionService) StartAuctionByPositions(roomID string, posMap map[string]int) error {
	var positions []PositionAuction
	for pos, count := range posMap {
		players, err := storage.FetchRandomPlayersByPosition(pos, count)
		if err != nil {
			return err
		}
		positions = append(positions, PositionAuction{
			Position: pos,
			Players:  players,
			Index:    0,
		})
	}
	a.StateMutex.Lock()
	state := &AuctionState{Positions: positions, CurrentPos: 0}
	a.State[roomID] = state
	a.StateMutex.Unlock()
	a.broadcastNextPlayer(roomID)
	return nil
}

func (a *AuctionService) broadcastNextPlayer(roomID string) {
	a.StateMutex.Lock()
	state := a.State[roomID]
	if state.CurrentPos >= len(state.Positions) {
		a.StateMutex.Unlock()
		return // Auction complete
	}
	posAuction := &state.Positions[state.CurrentPos]
	if posAuction.Index >= len(posAuction.Players) {
		state.CurrentPos++
		if state.CurrentPos < len(state.Positions) {
			posAuction = &state.Positions[state.CurrentPos]
			posAuction.Index = 0
		} else {
			a.StateMutex.Unlock()
			return // Auction complete
		}
	}
	player := posAuction.Players[posAuction.Index]
	state.CurrentBid = 0
	state.CurrentBidder = ""
	if state.Timer != nil {
		state.Timer.Stop()
	}
	state.Timer = time.AfterFunc(10*time.Second, func() {
		a.finishAuction(roomID)
	})
	a.StateMutex.Unlock()
	a.Broadcast(roomID, "auctionPlayer", map[string]interface{}{
		"position": posAuction.Position,
		"player":   player,
	})
}

func (a *AuctionService) PlaceBid(roomID, userID string, bid int) bool {
	a.StateMutex.Lock()
	defer a.StateMutex.Unlock()
	state := a.State[roomID]
	if bid > state.CurrentBid {
		state.CurrentBid = bid
		state.CurrentBidder = userID
		return true
	}
	return false
}

func (a *AuctionService) finishAuction(roomID string) {
	a.StateMutex.Lock()
	state := a.State[roomID]
	posAuction := &state.Positions[state.CurrentPos]
	player := posAuction.Players[posAuction.Index]
	winner := state.CurrentBidder
	bid := state.CurrentBid
	posAuction.Index++
	complete := state.CurrentPos >= len(state.Positions) && posAuction.Index >= len(posAuction.Players)
	a.StateMutex.Unlock()

	if winner != "" && a.Redis != nil {
		a.Redis.AddPlayerToTeam(roomID, winner, player)
	}
	a.Broadcast(roomID, "playerSold", map[string]interface{}{
		"position": posAuction.Position,
		"player":   player,
		"winner":   winner,
		"bid":      bid,
	})

	if !complete {
		a.broadcastNextPlayer(roomID)
	}
}
