package storage

import (
	"auction-house/internal/auction"
	"errors"
	"github.com/google/uuid"
	"sync"
)

type memoryRepository struct {
	mu       sync.RWMutex
	auctions map[uuid.UUID]*auction.Auction
}

func NewMemoryRepository() *memoryRepository {
	return &memoryRepository{
		auctions: make(map[uuid.UUID]*auction.Auction),
	}
}

func (r *memoryRepository) Create(a *auction.Auction) error {
	if a == nil {
		return errors.New("auction cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.auctions[a.ID()]; exists {
		return errors.New("auction already exists")
	}

	r.auctions[a.ID()] = a
	return nil
}

func (r *memoryRepository) Save(a *auction.Auction) error {
	if a == nil {
		return errors.New("auction cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.auctions[a.ID()]; !exists {
		return errors.New("auction with this id does not exist")
	}

	r.auctions[a.ID()] = a
	return nil
}
