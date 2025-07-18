package service

import (
	"context"

	"github.com/Creative-genius001/Stacklo/services/user/model"
)

type Service interface {
	GetUser(ctx context.Context, id string) (*model.User, error)
	UpdateUser(ctx context.Context, w model.User) error
}

type userService struct {
	repository Repository
}

func NewUserService(r Repository) Service {
	return &userService{r}
}

func (u *userService) GetUser(ctx context.Context, id string) (*model.User, error) {
	res, err := u.repository.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *userService) UpdateUser(ctx context.Context, w model.User) error {
	err := u.repository.UpdateUser(ctx, w)
	if err != nil {
		return err
	}
	return nil
}
