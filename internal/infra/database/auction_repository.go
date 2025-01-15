package database

import (
	"context"
	"fmt"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity"
	"fullcycle-auction_go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	AuctionMongo struct {
		Id          string                  `bson:"_id"`
		ProductName string                  `bson:"product_name"`
		Category    string                  `bson:"category"`
		Description string                  `bson:"description"`
		Condition   entity.ProductCondition `bson:"condition"`
		Status      entity.AuctionStatus    `bson:"status"`
		Timestamp   int64                   `bson:"timestamp"`
	}

	AuctionRepository struct {
		Collection      *mongo.Collection
		auctionDuration time.Duration
	}
)

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	repository := &AuctionRepository{
		Collection:      database.Collection("auctions"),
		auctionDuration: getAuctionDuration(),
	}
	go repository.StartAuctionCloser(context.Background())
	return repository
}

func (ar *AuctionRepository) CreateAuction(ctx context.Context, auction *entity.Auction) error {
	auctionMongo := &AuctionMongo{
		Id:          auction.Id,
		ProductName: auction.ProductName,
		Category:    auction.Category,
		Description: auction.Description,
		Condition:   auction.Condition,
		Status:      auction.Status,
		Timestamp:   auction.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	return nil
}

func (ar *AuctionRepository) FindAuctionById(
	ctx context.Context, id string) (*entity.Auction, error) {
	filter := bson.M{"_id": id}

	var auctionMongo AuctionMongo
	if err := ar.Collection.FindOne(ctx, filter).Decode(&auctionMongo); err != nil {
		logger.Error(fmt.Sprintf("Error trying to find auction by id = %s", id), err)
		return nil, internal_error.NewInternalServerError("Error trying to find auction by id")
	}

	return &entity.Auction{
		Id:          auctionMongo.Id,
		ProductName: auctionMongo.ProductName,
		Category:    auctionMongo.Category,
		Description: auctionMongo.Description,
		Condition:   auctionMongo.Condition,
		Status:      auctionMongo.Status,
		Timestamp:   time.Unix(auctionMongo.Timestamp, 0),
	}, nil
}

func (ar *AuctionRepository) FindAuctions(
	ctx context.Context,
	status entity.AuctionStatus,
	category string,
	productName string) ([]entity.Auction, error) {
	filter := bson.M{}

	if status != 0 {
		filter["status"] = status
	}

	if category != "" {
		filter["category"] = category
	}

	if productName != "" {
		filter["productName"] = primitive.Regex{Pattern: productName, Options: "i"}
	}

	cursor, err := ar.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error("Error finding auctions", err)
		return nil, internal_error.NewInternalServerError("Error finding auctions")
	}
	defer cursor.Close(ctx)

	var auctionsMongo []AuctionMongo
	if err := cursor.All(ctx, &auctionsMongo); err != nil {
		logger.Error("Error decoding auctions", err)
		return nil, internal_error.NewInternalServerError("Error decoding auctions")
	}

	var auctions []entity.Auction
	for _, auction := range auctionsMongo {
		auctions = append(auctions, entity.Auction{
			Id:          auction.Id,
			ProductName: auction.ProductName,
			Category:    auction.Category,
			Status:      auction.Status,
			Description: auction.Description,
			Condition:   auction.Condition,
			Timestamp:   time.Unix(auction.Timestamp, 0),
		})
	}

	return auctions, nil
}

func (ar *AuctionRepository) StartAuctionCloser(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ar.closeExpiredAuctions(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (ar *AuctionRepository) closeExpiredAuctions(ctx context.Context) {
	filter := bson.M{"status": entity.Active}
	cursor, err := ar.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error("Error finding open auctions", err)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var auction AuctionMongo
		if err = cursor.Decode(&auction); err != nil {
			logger.Error("Error decoding auction", err)
			continue
		}

		timeStamp := time.Unix(auction.Timestamp, 0)

		if time.Now().After(timeStamp.Add(ar.auctionDuration)) {
			auction.Status = entity.Completed
			update := bson.M{"$set": bson.M{"status": entity.Completed}}
			if _, err := ar.Collection.UpdateOne(ctx, bson.M{"_id": auction.Id}, update); err != nil {
				logger.Error("Error updating auction status to completed", err)
			}
		}
	}
}

func getAuctionDuration() time.Duration {
	auctionDuration := os.Getenv("AUCTION_DURATION")
	duration, err := time.ParseDuration(auctionDuration)
	if err != nil {
		return time.Second * 30
	}
	return duration
}
