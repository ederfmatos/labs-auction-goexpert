package usecase

import (
	"context"
	"fullcycle-auction_go/internal/entity"
	"fullcycle-auction_go/internal/internal_error"
)

func NewUserUseCase(userRepository entity.UserRepositoryInterface) UserUseCaseInterface {
	return &UserUseCase{
		userRepository,
	}
}

type UserUseCase struct {
	UserRepository entity.UserRepositoryInterface
}

type UserOutputDTO struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UserUseCaseInterface interface {
	FindUserById(
		ctx context.Context,
		id string) (*UserOutputDTO, *internal_error.InternalError)
}

func (u *UserUseCase) FindUserById(
	ctx context.Context, id string) (*UserOutputDTO, *internal_error.InternalError) {
	userEntity, err := u.UserRepository.FindUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &UserOutputDTO{
		Id:   userEntity.Id,
		Name: userEntity.Name,
	}, nil
}
