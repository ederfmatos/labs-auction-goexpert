package repository

import (
	"context"
	"fullcycle-auction_go/internal/entity"
)

type AuctionRepository interface {
	CreateAuction(ctx context.Context, auction *entity.Auction) error

	FindAuctions(ctx context.Context, status entity.AuctionStatus, category, productName string) ([]entity.Auction, error)

	FindAuctionById(ctx context.Context, id string) (*entity.Auction, error)
}
