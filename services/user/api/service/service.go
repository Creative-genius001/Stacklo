package service

import (
	"context"
	er "errors"

	errors "github.com/Creative-genius001/Stacklo/services/transaction/utils/error"
	"github.com/Creative-genius001/Stacklo/services/user/model"
	"github.com/Creative-genius001/Stacklo/services/user/types"
	"github.com/Creative-genius001/Stacklo/services/user/utils"
	"github.com/Creative-genius001/Stacklo/services/user/utils/logger"
)

type Service interface {
	GetUser(ctx context.Context, id string) (*model.User, error)
	UpdateUser(ctx context.Context, id string, data types.UpdateUser) error
}

type userService struct {
	repository Repository
}

func NewUserService(r Repository) Service {
	return &userService{r}
}

func (u *userService) GetUser(ctx context.Context, id string) (*model.User, error) {
	isValid := utils.IsValidUUID(id)
	if !isValid {
		logger.Logger.Error("user id is invalid and could not be parsed")
		return nil, errors.Wrap(errors.TypeInvalidInput, "user id could not be parsed", er.New("user id is not valid"))
	}
	res, err := u.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *userService) UpdateUser(ctx context.Context, id string, data types.UpdateUser) error {
	return u.repository.UpdateUser(ctx, id, data)
}
