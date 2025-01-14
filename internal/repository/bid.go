package repository

import (
	"context"
	"fullcycle-auction_go/internal/entity"
)

type BidRepository interface {
	CreateBid(ctx context.Context, bidEntities []entity.Bid) error

	FindBidByAuctionId(ctx context.Context, auctionId string) ([]entity.Bid, error)

	FindWinningBidByAuctionId(ctx context.Context, auctionId string) (*entity.Bid, error)
}
