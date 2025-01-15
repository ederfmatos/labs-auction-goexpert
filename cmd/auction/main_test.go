package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"fullcycle-auction_go/internal/entity"
	"fullcycle-auction_go/internal/usecase"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuctionFlow(t *testing.T) {
	ctx := context.Background()
	mongoContainer, endpoint := startMongoContainer(t, ctx)

	clientOptions := options.Client().ApplyURI("mongodb://" + endpoint)
	mongoClient, err := mongo.Connect(ctx, clientOptions)
	require.NoError(t, err)
	defer mongoContainer.Terminate(ctx)

	databaseConnection := mongoClient.Database("test")
	router := makeRouter(databaseConnection)

	auctionJSON := `{ "product_name": "Product Test", "category": "Category Test", "description": "Description Test", "condition": 1 }`
	req := httptest.NewRequest("POST", "/auction", bytes.NewBufferString(auctionJSON))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	req = httptest.NewRequest("GET", "/auction?status=0", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var auctions []usecase.AuctionOutputDTO
	err = json.Unmarshal(rec.Body.Bytes(), &auctions)
	require.NoError(t, err)
	require.NotEmpty(t, auctions)

	auctionId := auctions[0].Id

	time.Sleep(1 * time.Minute)

	req = httptest.NewRequest("GET", fmt.Sprintf("/auction/%s", auctionId), nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var auction entity.Auction
	err = json.Unmarshal(rec.Body.Bytes(), &auction)
	require.NoError(t, err)

	require.Equal(t, entity.Completed, auction.Status)
}

func startMongoContainer(t *testing.T, ctx context.Context) (*mongodb.MongoDBContainer, string) {
	mongoContainer, err := mongodb.Run(ctx, "mongo:6")
	require.NoError(t, err)

	endpoint, err := mongoContainer.Endpoint(context.Background(), "")
	require.NoError(t, err)
	return mongoContainer, endpoint
}
