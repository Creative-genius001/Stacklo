package auth

import (
	"context"

	"github.com/Creative-genius001/Stacklo/services/user/api/service"
	"github.com/Creative-genius001/Stacklo/services/user/model"
)

type Auth interface {
	CreateUser(ctx context.Context, w model.User) (*model.User, error)
	Login(ctx context.Context, w model.User) error
}

type authService struct {
	repository service.Repository
}

// CreateUser implements Auth.
func (a *authService) CreateUser(ctx context.Context, w model.User) (*model.User, error) {
	panic("unimplemented")
}

// Login implements Auth.
func (a *authService) Login(ctx context.Context, w model.User) error {
	panic("unimplemented")
}

func NewAuthService(r service.Repository) Auth {
	return &authService{r}
}
