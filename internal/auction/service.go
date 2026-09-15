package auction

import (
	"fmt"
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
func (s *Service) HandleCreate(productName string, startingPrice float64, minIncrement float64, buyoutPrice float64) (uuid.UUID, error) {

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
