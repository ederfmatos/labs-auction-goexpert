package usecase

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity"
	"fullcycle-auction_go/internal/repository"
	"os"
	"strconv"
	"time"
)

type (
	BidInputDTO struct {
		UserId    string  `json:"user_id"`
		AuctionId string  `json:"auction_id"`
		Amount    float64 `json:"amount"`
	}

	BidOutputDTO struct {
		Id        string    `json:"id"`
		UserId    string    `json:"user_id"`
		AuctionId string    `json:"auction_id"`
		Amount    float64   `json:"amount"`
		Timestamp time.Time `json:"timestamp" time_format:"2006-01-02 15:04:05"`
	}

	bidUseCase struct {
		BidRepository repository.BidRepository

		timer               *time.Timer
		maxBatchSize        int
		batchInsertInterval time.Duration
		bidChannel          chan entity.Bid
	}

	BidUseCase interface {
		CreateBid(ctx context.Context, bidInputDTO BidInputDTO) error

		FindWinningBidByAuctionId(ctx context.Context, auctionId string) (*BidOutputDTO, error)

		FindBidByAuctionId(ctx context.Context, auctionId string) ([]BidOutputDTO, error)
	}
)

func NewBidUseCase(bidRepository repository.BidRepository) BidUseCase {
	maxSizeInterval := getMaxBatchSizeInterval()
	maxBatchSize := getMaxBatchSize()

	bidUseCase := &bidUseCase{
		BidRepository:       bidRepository,
		maxBatchSize:        maxBatchSize,
		batchInsertInterval: maxSizeInterval,
		timer:               time.NewTimer(maxSizeInterval),
		bidChannel:          make(chan entity.Bid, maxBatchSize),
	}

	go bidUseCase.triggerCreateRoutine(context.Background())

	return bidUseCase
}

var bidBatch []entity.Bid

func (bu *bidUseCase) triggerCreateRoutine(ctx context.Context) {
	defer close(bu.bidChannel)

	for {
		select {
		case bid, ok := <-bu.bidChannel:
			if !ok {
				if len(bidBatch) > 0 {
					if err := bu.BidRepository.CreateBid(ctx, bidBatch); err != nil {
						logger.Error("error trying to process bid batch list", err)
					}
				}
				return
			}

			bidBatch = append(bidBatch, bid)

			if len(bidBatch) >= bu.maxBatchSize {
				if err := bu.BidRepository.CreateBid(ctx, bidBatch); err != nil {
					logger.Error("error trying to process bid batch list", err)
				}

				bidBatch = nil
				bu.timer.Reset(bu.batchInsertInterval)
			}
		case <-bu.timer.C:
			if err := bu.BidRepository.CreateBid(ctx, bidBatch); err != nil {
				logger.Error("error trying to process bid batch list", err)
			}
			bidBatch = nil
			bu.timer.Reset(bu.batchInsertInterval)
		}
	}
}

func (bu *bidUseCase) CreateBid(ctx context.Context, input BidInputDTO) error {
	bid, err := entity.CreateBid(input.UserId, input.AuctionId, input.Amount)
	if err != nil {
		return err
	}

	bu.bidChannel <- *bid

	return nil
}

func getMaxBatchSizeInterval() time.Duration {
	batchInsertInterval := os.Getenv("BATCH_INSERT_INTERVAL")
	duration, err := time.ParseDuration(batchInsertInterval)
	if err != nil {
		return 3 * time.Minute
	}
	return duration
}

func getMaxBatchSize() int {
	value, err := strconv.Atoi(os.Getenv("MAX_BATCH_SIZE"))
	if err != nil {
		return 5
	}
	return value
}
