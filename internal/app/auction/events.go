package auction

import (
	"github.com/yourusername/TouchlineTactics/internal/storage"
)

type StartAuctionPayload struct {
	RoomID     string `json:"roomId"`
	NumPlayers int    `json:"numPlayers"`
}

type PlaceBidPayload struct {
	RoomID string `json:"roomId"`
	UserID string `json:"userId"`
	Bid    int    `json:"bid"`
}

type AuctionEventHandler struct {
	Auction *AuctionService
}

func (h *AuctionEventHandler) HandleStartAuction(payload StartAuctionPayload) error {
	players, err := storage.FetchRandomPlayers(payload.NumPlayers)
	if err != nil {
		return err
	}
	if err2 := h.Auction.Start(payload.RoomID, players); err2 != nil {
		if e, ok := err2.(error); ok {
			return e
		}
		return nil
	}
	return nil
}

func (h *AuctionEventHandler) HandlePlaceBid(payload PlaceBidPayload) {
	h.Auction.PlaceBid(payload.RoomID, payload.UserID, payload.Bid)
}
