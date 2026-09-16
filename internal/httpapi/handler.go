package httpapi

import (
	"auction-house/internal/auction"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Handler struct {
	service *auction.Service
}

func NewHandler(service *auction.Service) *Handler {
	return &Handler{service: service}
}

type createRequest struct {
	ProductName   string `json:"productName"`
	StartingPrice int64  `json:"startingPrice"`
	MinIncrement  int64  `json:"minIncrement"`
	BuyoutPrice   int64  `json:"buyoutPrice"`
}

type createResponse struct {
	AuctionID uuid.UUID `json:"auctionId"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {

	var request createRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if request.ProductName == "" || request.StartingPrice <= 0 || request.BuyoutPrice <= 0 || request.MinIncrement <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "productName, startingPrice, buyoutPrice and minIncrement are required"})
		return
	}

	id, err := h.service.HandleCreate(request.ProductName, auction.Money(request.StartingPrice), auction.Money(request.MinIncrement), auction.Money(request.BuyoutPrice))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not create auction"})
		return
	}

	writeJSON(w, http.StatusCreated, createResponse{AuctionID: id})
}

type getBidResponse struct {
	BidID      uuid.UUID `json:"bidId"`
	Amount     int64     `json:"amount"`
	AcceptedAt time.Time `json:"acceptedAt"`
}

type getResponse struct {
	AuctionID     uuid.UUID        `json:"auctionId"`
	ProductName   string           `json:"productName"`
	StartingPrice int64            `json:"startingPrice"`
	MinIncrement  int64            `json:"minIncrement"`
	BuyoutPrice   int64            `json:"buyoutPrice"`
	Opened        bool             `json:"opened"`
	Bids          []getBidResponse `json:"bids"`
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {

	rawID := r.PathValue("id")

	id, err := uuid.Parse(rawID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid auction id"})
		return
	}

	auction, err := h.service.HandleGet(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	domainBids := auction.Bids()

	bids := make([]getBidResponse, 0, len(domainBids))

	for _, bid := range domainBids {
		bids = append(bids, getBidResponse{
			BidID:      bid.ID(),
			Amount:     int64(bid.Amount()),
			AcceptedAt: bid.AcceptedAt(),
		})
	}

	response := getResponse{
		AuctionID:     auction.ID(),
		ProductName:   auction.ProductName(),
		StartingPrice: int64(auction.StartingPrice()),
		MinIncrement:  int64(auction.MinIncrement()),
		BuyoutPrice:   int64(auction.BuyoutPrice()),
		Opened:        auction.Opened(),
		Bids:          bids,
	}

	writeJSON(w, http.StatusOK, response)
}

type bidRequest struct {
	Amount int64 `json:"amount"`
}

type bidIDResponse struct {
	BidID uuid.UUID `json:"bidId"`
}

func (h *Handler) bid(w http.ResponseWriter, r *http.Request) {

	var request bidRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	rawID := r.PathValue("id")
	id, err := uuid.Parse(rawID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid auction id"})
		return
	}

	if request.Amount <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "amount must be positive"})
		return
	}

	bidId, err := h.service.HandleBid(id, auction.Money(request.Amount))
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, bidIDResponse{BidID: bidId})
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /auctions", h.create)
	mux.HandleFunc("POST /auctions/{id}/bids", h.bid)
	mux.HandleFunc("GET /auctions/{id}", h.get)

	return mux
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, auction.ErrAuctionNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "auction not found"})
	case errors.Is(err, auction.ErrBidTooLow):
		writeJSON(w, http.StatusConflict, map[string]string{"error": "bid amount is too low"})
	case errors.Is(err, auction.ErrAuctionClosed):
		writeJSON(w, http.StatusConflict, map[string]string{"error": "auction is closed"})
	case errors.Is(err, auction.ErrDuplicateBid):
		writeJSON(w, http.StatusConflict, map[string]string{"error": "bid already exists"})
	default:
		log.Printf("unexpected service error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
}
