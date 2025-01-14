package repository

import (
	"context"
	"fullcycle-auction_go/internal/entity"
)

type UserRepository interface {
	FindUserById(ctx context.Context, userId string) (*entity.User, error)
}
