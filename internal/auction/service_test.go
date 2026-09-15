package auction_test

import (
	"auction-house/internal/auction"
	"auction-house/internal/storage"
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestAuctionBidFlow(t *testing.T) {
	repo := storage.NewMemoryRepository()
	service := auction.NewService(repo)

	auctionID, err := service.HandleCreate(
		"ferarri testarossa",
		auction.Money(100_000_00), // starting price 100k
		auction.Money(5_000_00),   // minimum increment 5k
		auction.Money(200_000_00), // buyout price 200k
	)
	if err != nil {
		t.Fatalf("create auction: %v", err)
	}

	if auctionID == uuid.Nil {
		t.Fatalf("get created auction: %v", err)
	}

	storedAuction, err := service.HandleGet(auctionID)
	if err != nil {
		t.Fatalf("get created auction: %v", err)
	}

	if storedAuction.ID() != auctionID {
		t.Fatalf("expected auction ID %s, got %s", auctionID, storedAuction.ID())
	}

	// A first bid below the starting price must fail
	_, err = service.HandleBid(auctionID, auction.Money(90_000_00)) // 90k
	if !errors.Is(err, auction.ErrBidTooLow) {
		t.Fatalf("expected ErrBidTooLow, got %v", err)
	}

	// A first bid equal to the starting price is valid
	bidID, err := service.HandleBid(auctionID, auction.Money(100_000_00)) // 100k
	if err != nil {
		t.Fatalf("place first bid: %v", err)
	}

	if bidID == uuid.Nil {
		t.Fatalf("expected a non-empty bid ID")
	}

	// The next bid must be at least 105k
	_, err = service.HandleBid(auctionID, auction.Money(104_999_99)) // 104 999.99
	if !errors.Is(err, auction.ErrBidTooLow) {
		t.Fatalf("expected ErrBidTooLow, got %v", err)
	}

	// Reaching the buyout price accepts the bid and closes the auction
	_, err = service.HandleBid(auctionID, auction.Money(200_000_00)) // 200k
	if err != nil {
		t.Fatalf("place buyout bid: %v", err)
	}

	// further bids mut be rejected
	_, err = service.HandleBid(auctionID, auction.Money(205_000_00)) // 205k
	if !errors.Is(err, auction.ErrAuctionClosed) {
		t.Fatalf("expected ErrAuctionClosed, got %v", err)
	}
}
