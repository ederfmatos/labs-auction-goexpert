package entity

import (
	"fullcycle-auction_go/internal/internal_error"
	"github.com/google/uuid"
	"time"
)

type Bid struct {
	Id        string
	UserId    string
	AuctionId string
	Amount    float64
	Timestamp time.Time
}

func CreateBid(userId, auctionId string, amount float64) (*Bid, error) {
	bid := &Bid{
		Id:        uuid.New().String(),
		UserId:    userId,
		AuctionId: auctionId,
		Amount:    amount,
		Timestamp: time.Now(),
	}

	if err := bid.Validate(); err != nil {
		return nil, err
	}

	return bid, nil
}

func (b *Bid) Validate() error {
	if err := uuid.Validate(b.UserId); err != nil {
		return internal_error.NewBadRequestError("UserId is not a valid id")
	}
	if err := uuid.Validate(b.AuctionId); err != nil {
		return internal_error.NewBadRequestError("AuctionId is not a valid id")
	}
	if b.Amount <= 0 {
		return internal_error.NewBadRequestError("Amount is not a valid value")
	}
	return nil
}
