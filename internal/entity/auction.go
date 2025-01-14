package entity

import (
	"fullcycle-auction_go/internal/internal_error"
	"github.com/google/uuid"
	"time"
)

const (
	Active AuctionStatus = iota
	Completed
)

const (
	New ProductCondition = iota + 1
	Used
	Refurbished
)

type (
	Auction struct {
		Id          string
		ProductName string
		Category    string
		Description string
		Condition   ProductCondition
		Status      AuctionStatus
		Timestamp   time.Time
	}

	ProductCondition int
	AuctionStatus    int
)

func CreateAuction(productName, category, description string, condition ProductCondition) (*Auction, error) {
	auction := &Auction{
		Id:          uuid.New().String(),
		ProductName: productName,
		Category:    category,
		Description: description,
		Condition:   condition,
		Status:      Active,
		Timestamp:   time.Now(),
	}

	if err := auction.Validate(); err != nil {
		return nil, err
	}

	return auction, nil
}

func (au *Auction) Validate() error {
	if len(au.ProductName) <= 1 ||
		len(au.Category) <= 2 ||
		len(au.Description) <= 10 && (au.Condition != New && au.Condition != Refurbished && au.Condition != Used) {
		return internal_error.NewBadRequestError("invalid auction object")
	}

	return nil
}
