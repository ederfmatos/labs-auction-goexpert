package usecase

import (
	"context"
	"fullcycle-auction_go/internal/entity"
	"fullcycle-auction_go/internal/repository"
	"time"
)

type (
	AuctionInputDTO struct {
		ProductName string           `json:"product_name" binding:"required,min=1"`
		Category    string           `json:"category" binding:"required,min=2"`
		Description string           `json:"description" binding:"required,min=10,max=200"`
		Condition   ProductCondition `json:"condition" binding:"oneof=0 1 2"`
	}

	AuctionOutputDTO struct {
		Id          string           `json:"id"`
		ProductName string           `json:"product_name"`
		Category    string           `json:"category"`
		Description string           `json:"description"`
		Condition   ProductCondition `json:"condition"`
		Status      AuctionStatus    `json:"status"`
		Timestamp   time.Time        `json:"timestamp" time_format:"2006-01-02 15:04:05"`
	}

	WinningInfoOutputDTO struct {
		Auction AuctionOutputDTO `json:"auction"`
		Bid     *BidOutputDTO    `json:"bid,omitempty"`
	}

	AuctionUseCase interface {
		CreateAuction(ctx context.Context, auctionInput AuctionInputDTO) error

		FindAuctionById(ctx context.Context, id string) (*AuctionOutputDTO, error)

		FindAuctions(ctx context.Context, status AuctionStatus, category, productName string) ([]AuctionOutputDTO, error)

		FindWinningBidByAuctionId(ctx context.Context, auctionId string) (*WinningInfoOutputDTO, error)
	}

	ProductCondition int64
	AuctionStatus    int64

	auctionUseCase struct {
		auctionRepository repository.AuctionRepository
		bidRepository     repository.BidRepository
	}
)

func NewAuctionUseCase(
	auctionRepository repository.AuctionRepository,
	bidRepository repository.BidRepository,
) AuctionUseCase {
	return &auctionUseCase{
		auctionRepository: auctionRepository,
		bidRepository:     bidRepository,
	}
}

func (au *auctionUseCase) CreateAuction(ctx context.Context, input AuctionInputDTO) error {
	auction, err := entity.CreateAuction(
		input.ProductName,
		input.Category,
		input.Description,
		entity.ProductCondition(input.Condition),
	)
	if err != nil {
		return err
	}

	if err := au.auctionRepository.CreateAuction(ctx, auction); err != nil {
		return err
	}
	return nil
}
