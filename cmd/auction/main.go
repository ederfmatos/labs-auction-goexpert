package main

import (
	"fullcycle-auction_go/configuration/database/mongodb"
	"fullcycle-auction_go/internal/infra/api"
	"fullcycle-auction_go/internal/infra/database"
	"fullcycle-auction_go/internal/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func main() {
	databaseConnection, err := mongodb.NewMongoDBConnection()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	router := makeRouter(databaseConnection)
	if err = router.Run(":8080"); err != nil {
		log.Fatalf("Error trying to start server: %s", err.Error())
	}
}

func makeRouter(databaseConnection *mongo.Database) *gin.Engine {
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

	return router
}
