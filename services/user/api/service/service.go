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

// Close implements Service.
func (u *userService) Close() {
	panic("unimplemented")
}

func (u *userService) CreateUser(ctx context.Context, w model.User) (*model.User, error) {
	panic("unimplemented")
}

// GetUser implements Service.
func (u *userService) GetUser(ctx context.Context, id string) (*model.User, error) {
	panic("unimplemented")
}

// UpdateUser implements Service.
func (u *userService) UpdateUser(ctx context.Context, w model.User) error {
	panic("unimplemented")
}

func NewUserService(r Repository) Service {
	return &userService{r}
}
