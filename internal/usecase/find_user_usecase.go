package usecase

import (
	"context"
	"fullcycle-auction_go/internal/repository"
)

type (
	userUseCase struct {
		userRepository repository.UserRepository
	}

	UserOutputDTO struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}

	UserUseCase interface {
		FindUserById(ctx context.Context, id string) (*UserOutputDTO, error)
	}
)

func NewUserUseCase(userRepository repository.UserRepository) UserUseCase {
	return &userUseCase{userRepository: userRepository}
}

func (u *userUseCase) FindUserById(ctx context.Context, id string) (*UserOutputDTO, error) {
	user, err := u.userRepository.FindUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &UserOutputDTO{
		Id:   user.Id,
		Name: user.Name,
	}, nil
}
