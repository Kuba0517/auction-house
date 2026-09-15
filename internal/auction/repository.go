package auction

import "github.com/google/uuid"

type Repository interface {
	Create(*Auction) error
	Save(*Auction) error
	Get(uuid.UUID) (*Auction, error)
}
