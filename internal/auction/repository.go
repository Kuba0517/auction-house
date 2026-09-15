package auction

type Repository interface {
	Create(*Auction) error
	Save(*Auction) error
}
