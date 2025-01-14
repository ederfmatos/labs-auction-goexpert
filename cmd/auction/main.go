package main

import (
	"context"
	"fullcycle-auction_go/configuration/database/mongodb"
	"fullcycle-auction_go/internal/infra/api"
	"fullcycle-auction_go/internal/infra/database"
	"fullcycle-auction_go/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load("cmd/auction/.env"); err != nil {
		log.Fatal("Error trying to load env variables")
		return
	}

	databaseConnection, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	router := gin.Default()

	auctionRepository := database.NewAuctionRepository(databaseConnection)
	bidRepository := database.NewBidRepository(databaseConnection, auctionRepository)
	userRepository := database.NewUserRepository(databaseConnection)

	userController := api.NewUserController(usecase.NewUserUseCase(userRepository))
	auctionController := api.NewAuctionController(usecase.NewAuctionUseCase(auctionRepository, bidRepository))
	bidController := api.NewBidController(usecase.NewBidUseCase(bidRepository))

	router.GET("/auction", auctionController.FindAuctions)
	router.GET("/auction/:auctionId", auctionController.FindAuctionById)
	router.POST("/auction", auctionController.CreateAuction)
	router.GET("/auction/winner/:auctionId", auctionController.FindWinningBidByAuctionId)
	router.POST("/bid", bidController.CreateBid)
	router.GET("/bid/:auctionId", bidController.FindBidByAuctionId)
	router.GET("/user/:userId", userController.FindUserById)

	router.Run(":8080")
}
