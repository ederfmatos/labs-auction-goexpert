package database

import (
	"context"
	"fmt"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity"
	"fullcycle-auction_go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	BidMongo struct {
		Id        string  `bson:"_id"`
		UserId    string  `bson:"user_id"`
		AuctionId string  `bson:"auction_id"`
		Amount    float64 `bson:"amount"`
		Timestamp int64   `bson:"timestamp"`
	}

	BidRepository struct {
		Collection            *mongo.Collection
		AuctionRepository     *AuctionRepository
		auctionInterval       time.Duration
		auctionStatusMap      map[string]entity.AuctionStatus
		auctionEndTimeMap     map[string]time.Time
		auctionStatusMapMutex *sync.Mutex
		auctionEndTimeMutex   *sync.Mutex
	}
)

func NewBidRepository(database *mongo.Database, auctionRepository *AuctionRepository) *BidRepository {
	return &BidRepository{
		auctionInterval:       getAuctionInterval(),
		auctionStatusMap:      make(map[string]entity.AuctionStatus),
		auctionEndTimeMap:     make(map[string]time.Time),
		auctionStatusMapMutex: &sync.Mutex{},
		auctionEndTimeMutex:   &sync.Mutex{},
		Collection:            database.Collection("bids"),
		AuctionRepository:     auctionRepository,
	}
}

func (bd *BidRepository) CreateBid(ctx context.Context, bidEntities []entity.Bid) error {
	var wg sync.WaitGroup
	for _, bid := range bidEntities {
		wg.Add(1)
		go func(bidValue entity.Bid) {
			defer wg.Done()

			bd.auctionStatusMapMutex.Lock()
			auctionStatus, okStatus := bd.auctionStatusMap[bidValue.AuctionId]
			bd.auctionStatusMapMutex.Unlock()

			bd.auctionEndTimeMutex.Lock()
			auctionEndTime, okEndTime := bd.auctionEndTimeMap[bidValue.AuctionId]
			bd.auctionEndTimeMutex.Unlock()

			bidMongo := &BidMongo{
				Id:        bidValue.Id,
				UserId:    bidValue.UserId,
				AuctionId: bidValue.AuctionId,
				Amount:    bidValue.Amount,
				Timestamp: bidValue.Timestamp.Unix(),
			}

			if okEndTime && okStatus {
				now := time.Now()
				if auctionStatus == entity.Completed || now.After(auctionEndTime) {
					return
				}

				if _, err := bd.Collection.InsertOne(ctx, bidMongo); err != nil {
					logger.Error("Error trying to insert bid", err)
					return
				}

				return
			}

			auction, err := bd.AuctionRepository.FindAuctionById(ctx, bidValue.AuctionId)
			if err != nil {
				logger.Error("Error trying to find auction by id", err)
				return
			}
			if auction.Status == entity.Completed {
				return
			}

			bd.auctionStatusMapMutex.Lock()
			bd.auctionStatusMap[bidValue.AuctionId] = auction.Status
			bd.auctionStatusMapMutex.Unlock()

			bd.auctionEndTimeMutex.Lock()
			bd.auctionEndTimeMap[bidValue.AuctionId] = auction.Timestamp.Add(bd.auctionInterval)
			bd.auctionEndTimeMutex.Unlock()

			if _, err := bd.Collection.InsertOne(ctx, bidMongo); err != nil {
				logger.Error("Error trying to insert bid", err)
				return
			}
		}(bid)
	}
	wg.Wait()
	return nil
}

func (bd *BidRepository) FindBidByAuctionId(
	ctx context.Context, auctionId string) ([]entity.Bid, error) {
	filter := bson.M{"auctionId": auctionId}

	cursor, err := bd.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying to find bids by auctionId %s", auctionId), err)
		return nil, internal_error.NewInternalServerError(fmt.Sprintf("Error trying to find bids by auctionId %s", auctionId))
	}

	var bids []BidMongo
	if err := cursor.All(ctx, &bids); err != nil {
		logger.Error(fmt.Sprintf("Error trying to find bids by auctionId %s", auctionId), err)
		return nil, internal_error.NewInternalServerError(fmt.Sprintf("Error trying to find bids by auctionId %s", auctionId))
	}

	var bidEntities []entity.Bid
	for _, bidMongo := range bids {
		bidEntities = append(bidEntities, entity.Bid{
			Id:        bidMongo.Id,
			UserId:    bidMongo.UserId,
			AuctionId: bidMongo.AuctionId,
			Amount:    bidMongo.Amount,
			Timestamp: time.Unix(bidMongo.Timestamp, 0),
		})
	}

	return bidEntities, nil
}

func (bd *BidRepository) FindWinningBidByAuctionId(
	ctx context.Context, auctionId string) (*entity.Bid, error) {
	filter := bson.M{"auction_id": auctionId}

	var bidMongo BidMongo
	opts := options.FindOne().SetSort(bson.D{{"amount", -1}})
	if err := bd.Collection.FindOne(ctx, filter, opts).Decode(&bidMongo); err != nil {
		logger.Error("Error trying to find the auction winner", err)
		return nil, internal_error.NewInternalServerError("Error trying to find the auction winner")
	}

	return &entity.Bid{
		Id:        bidMongo.Id,
		UserId:    bidMongo.UserId,
		AuctionId: bidMongo.AuctionId,
		Amount:    bidMongo.Amount,
		Timestamp: time.Unix(bidMongo.Timestamp, 0),
	}, nil
}

func getAuctionInterval() time.Duration {
	auctionInterval := os.Getenv("AUCTION_INTERVAL")
	duration, err := time.ParseDuration(auctionInterval)
	if err != nil {
		return time.Minute * 5
	}

	return duration
}
