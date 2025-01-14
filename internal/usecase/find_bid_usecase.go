package usecase

import (
	"context"
)

func (bu *bidUseCase) FindBidByAuctionId(ctx context.Context, auctionId string) ([]BidOutputDTO, error) {
	bidList, err := bu.BidRepository.FindBidByAuctionId(ctx, auctionId)
	if err != nil {
		return nil, err
	}

	var bidOutputList []BidOutputDTO
	for _, bid := range bidList {
		bidOutputList = append(bidOutputList, BidOutputDTO{
			Id:        bid.Id,
			UserId:    bid.UserId,
			AuctionId: bid.AuctionId,
			Amount:    bid.Amount,
			Timestamp: bid.Timestamp,
		})
	}

	return bidOutputList, nil
}

func (bu *bidUseCase) FindWinningBidByAuctionId(
	ctx context.Context, auctionId string) (*BidOutputDTO, error) {
	bid, err := bu.BidRepository.FindWinningBidByAuctionId(ctx, auctionId)
	if err != nil {
		return nil, err
	}

	bidOutput := &BidOutputDTO{
		Id:        bid.Id,
		UserId:    bid.UserId,
		AuctionId: bid.AuctionId,
		Amount:    bid.Amount,
		Timestamp: bid.Timestamp,
	}

	return bidOutput, nil
}
