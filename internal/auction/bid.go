package auction

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type bid struct {
	id                uuid.UUID
	amount            float64
	acceptedTimestamp time.Time
}

func newBid(id uuid.UUID, amount float64, acceptedTimestamp time.Time) (*bid, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("id cannot be empty")
	}

	if amount <= 0 {
		return nil, fmt.Errorf("amount cannot be lower or equal to 0")
	}

	return &bid{
		id:                id,
		amount:            amount,
		acceptedTimestamp: acceptedTimestamp,
	}, nil
}
