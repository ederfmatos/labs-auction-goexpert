package usecase

import (
	"context"
	"fullcycle-auction_go/internal/entity"
	"time"
)

type AuctionInputDTO struct {
	ProductName string           `json:"product_name" binding:"required,min=1"`
	Category    string           `json:"category" binding:"required,min=2"`
	Description string           `json:"description" binding:"required,min=10,max=200"`
	Condition   ProductCondition `json:"condition" binding:"oneof=0 1 2"`
}

type AuctionOutputDTO struct {
	Id          string           `json:"id"`
	ProductName string           `json:"product_name"`
	Category    string           `json:"category"`
	Description string           `json:"description"`
	Condition   ProductCondition `json:"condition"`
	Status      AuctionStatus    `json:"status"`
	Timestamp   time.Time        `json:"timestamp" time_format:"2006-01-02 15:04:05"`
}

type WinningInfoOutputDTO struct {
	Auction AuctionOutputDTO `json:"auction"`
	Bid     *BidOutputDTO    `json:"bid,omitempty"`
}

func NewAuctionUseCase(
	auctionRepositoryInterface entity.AuctionRepositoryInterface,
	bidRepositoryInterface entity.BidEntityRepository) AuctionUseCaseInterface {
	return &AuctionUseCase{
		auctionRepositoryInterface: auctionRepositoryInterface,
		bidRepositoryInterface:     bidRepositoryInterface,
	}
}

type AuctionUseCaseInterface interface {
	CreateAuction(
		ctx context.Context,
		auctionInput AuctionInputDTO) error

	FindAuctionById(
		ctx context.Context, id string) (*AuctionOutputDTO, error)

	FindAuctions(
		ctx context.Context,
		status AuctionStatus,
		category, productName string) ([]AuctionOutputDTO, error)

	FindWinningBidByAuctionId(
		ctx context.Context,
		auctionId string) (*WinningInfoOutputDTO, error)
}

type ProductCondition int64
type AuctionStatus int64

type AuctionUseCase struct {
	auctionRepositoryInterface entity.AuctionRepositoryInterface
	bidRepositoryInterface     entity.BidEntityRepository
}

func (au *AuctionUseCase) CreateAuction(
	ctx context.Context,
	auctionInput AuctionInputDTO) error {
	auction, err := entity.CreateAuction(
		auctionInput.ProductName,
		auctionInput.Category,
		auctionInput.Description,
		entity.ProductCondition(auctionInput.Condition))
	if err != nil {
		return err
	}

	if err := au.auctionRepositoryInterface.CreateAuction(
		ctx, auction); err != nil {
		return err
	}

	return nil
}
