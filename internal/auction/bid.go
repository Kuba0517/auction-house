package auction

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Bid struct {
	id         uuid.UUID
	amount     Money
	acceptedAt time.Time
}

func newBid(id uuid.UUID, amount Money, acceptedAt time.Time) (*Bid, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("id cannot be empty")
	}

	if amount <= 0 {
		return nil, fmt.Errorf("amount cannot be lower or equal to 0")
	}

	return &Bid{
		id:         id,
		amount:     amount,
		acceptedAt: acceptedAt,
	}, nil
}

func (b Bid) ID() uuid.UUID {
	return b.id
}

func (b Bid) Amount() Money {
	return b.amount
}

func (b Bid) AcceptedAt() time.Time {
	return b.acceptedAt
}
