package database

import (
	"context"
	"errors"
	"fmt"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity"
	"fullcycle-auction_go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	UserMongo struct {
		Id   string `bson:"_id"`
		Name string `bson:"name"`
	}

	UserRepository struct {
		Collection *mongo.Collection
	}
)

func NewUserRepository(database *mongo.Database) *UserRepository {
	return &UserRepository{
		Collection: database.Collection("users"),
	}
}

func (ur *UserRepository) FindUserById(ctx context.Context, userId string) (*entity.User, error) {
	filter := bson.M{"_id": userId}

	var userMongo UserMongo
	err := ur.Collection.FindOne(ctx, filter).Decode(&userMongo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Error(fmt.Sprintf("User not found with this id = %d", userId), err)
			return nil, internal_error.NewNotFoundError(
				fmt.Sprintf("User not found with this id = %d", userId))
		}

		logger.Error("Error trying to find user by userId", err)
		return nil, internal_error.NewInternalServerError("Error trying to find user by userId")
	}

	user := &entity.User{
		Id:   userMongo.Id,
		Name: userMongo.Name,
	}

	return user, nil
}
