package auction

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Auction struct {
	id            uuid.UUID
	productName   string
	minIncrement  Money
	startingPrice Money
	buyoutPrice   Money
	opened        bool
	bids          []Bid
}

type Money int64

func newAuction(id uuid.UUID, productName string, minIncrement Money, startingPrice Money, buyoutPrice Money) (*Auction, error) {
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
		bids:          []Bid{},
	}, nil
}

func (a *Auction) PlaceBid(bid *Bid) error {

	if !a.opened {
		return ErrAuctionClosed
	}

	minimumAmount := a.startingPrice
	if len(a.bids) > 0 {
		previousBid := a.bids[len(a.bids)-1].amount
		minimumAmount = previousBid + a.minIncrement
	}

	if bid.amount < minimumAmount {
		return fmt.Errorf("%w: minimum acceptable amount is %s", ErrBidTooLow, minimumAmount)
	}

	a.bids = append(a.bids, *bid)

	if bid.amount >= a.buyoutPrice {
		a.opened = false
	}

	return nil
}

func (m Money) String() string {
	cents := Money(m)

	return fmt.Sprintf(
		"%d.%02d",
		cents/100,
		cents%100,
	)
}

func (a Auction) ID() uuid.UUID {
	return a.id
}

func (a Auction) String() string {
	return fmt.Sprintf(
		"Auction{ID: %s, Product: %q, StartingPrice: %s, MinIncrement: %s, BuyoutPrice: %s, Opened: %t, Bids: %d}",
		a.id.String(),
		a.productName,
		a.startingPrice.String(),
		a.minIncrement.String(),
		a.buyoutPrice.String(),
		a.opened,
		len(a.bids),
	)
}
