package auction

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// tries to create auction struct and save it into data source of Repository
func (s *Service) HandleCreate(productName string, startingPrice Money, minIncrement Money, buyoutPrice Money) (uuid.UUID, error) {

	auctionId := uuid.New()
	auction, err := newAuction(auctionId, productName, minIncrement, startingPrice, buyoutPrice)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create auction: %w", err)
	}

	if err := s.repo.Create(auction); err != nil {
		return uuid.Nil, fmt.Errorf("failed to store auction: %w", err)
	}

	return auctionId, nil
}

func (s *Service) HandleGet(id uuid.UUID) (*Auction, error) {
	return s.repo.Get(id)
}

func (s *Service) HandleBid(auctiondId uuid.UUID, amount Money) (uuid.UUID, error) {

	auction, err := s.repo.Get(auctiondId)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get auction: %w", err)
	}

	bidId := uuid.New()
	bid, err := newBid(bidId, amount, time.Now())
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create a bid: %w", err)
	}

	if err := auction.PlaceBid(bid); err != nil {
		return uuid.Nil, fmt.Errorf("failed to place bid: %w", err)
	}

	if err := s.repo.Save(auction); err != nil {
		return uuid.Nil, fmt.Errorf("failed to save auction: %w", err)
	}

	return bidId, nil
}
