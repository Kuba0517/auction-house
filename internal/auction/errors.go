package auction

import "errors"

var (
	ErrAuctionClosed   = errors.New("auction is closed")
	ErrAuctionNotFound = errors.New("auction was not found")
	ErrBidTooLow       = errors.New("bid amount is too low")
	ErrDuplicateBid    = errors.New("bid already exists")
)
