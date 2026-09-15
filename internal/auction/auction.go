package auction

import (
	"errors"
	"github.com/google/uuid"
)

type Auction struct {
	id            uuid.UUID
	productName   string
	minIncrement  float64
	startingPrice float64
	buyoutPrice   float64
	opened        bool
	bids          []bid
}

func newAuction(id uuid.UUID, productName string, minIncrement float64, startingPrice float64, buyoutPrice float64) (*Auction, error) {
	if id == uuid.Nil {
		return nil, errors.New("id cannot be empty")
	}

	if productName == "" {
		return nil, errors.New("productName cannot be empty")
	}

	if minIncrement <= 0 {
		return nil, errors.New("minIncrement cannot be lower or equal to 0")
	}

	if startingPrice <= 0 {
		return nil, errors.New("startingPrice cannot be lower or equal to 0")
	}

	if buyoutPrice <= 0 {
		return nil, errors.New("buyoutPrice cannot be lower or equal to 0")
	}

	return &Auction{
		id:            id,
		productName:   productName,
		minIncrement:  minIncrement,
		startingPrice: startingPrice,
		buyoutPrice:   buyoutPrice,
		opened:        true,
		bids:          []bid{},
	}, nil
}

func (a Auction) ID() uuid.UUID {
	return a.id
}
