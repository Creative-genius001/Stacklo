package auth

import (
	"context"
	er "errors"

	"github.com/Creative-genius001/Stacklo/services/user/api/service"
	"github.com/Creative-genius001/Stacklo/services/user/model"
	"github.com/Creative-genius001/Stacklo/services/user/utils"
	errors "github.com/Creative-genius001/Stacklo/services/user/utils/error"
	"github.com/Creative-genius001/Stacklo/services/user/utils/logger"
)

type Auth interface {
	Register(ctx context.Context, user model.User) (*model.User, error)
	Login(ctx context.Context, w model.User) error
}

type authService struct {
	repository service.Repository
}

func NewAuthService(r service.Repository) Auth {
	return &authService{r}
}

func (a *authService) Register(ctx context.Context, user model.User) (*model.User, error) {

	//validate email
	isValid := utils.IsValidEmail(user.Email)
	if !isValid {
		logger.Logger.Warn("invalid email address")
		return nil, errors.Wrap(errors.TypeInvalidInput, "Email address is invalid", er.New("Email address is invalid"))
	}

	//validate phoneNumber
	isValid, formatPhone, err := utils.IsValidPhoneNumber(user.PhoneNumber, "NG")
	if !isValid || err != nil {
		logger.Logger.Warn("invalid phone number")
		return nil, errors.Wrap(errors.TypeInvalidInput, "Phone number is invalid", er.New("Phone number is invalid"))
	}
	user.PhoneNumber = formatPhone

	//validate password
	isValidPassword := utils.IsValidPassword(user.PasswordHash)
	if !isValidPassword {
		logger.Logger.Warn("invalid password")
		return nil, errors.Wrap(errors.TypeInvalidInput, "invalid password", er.New("invalid passwordd"))
	}

	res, err := a.repository.FindByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	if res != nil {
		logger.Logger.Warn("user with this email already exists")
		return nil, errors.Wrap(errors.TypeConflict, "user with this email already exists", er.New("user with this email already exists"))
	}

	hashedPassword, _ := utils.HashPassword(user.PasswordHash)
	user.PasswordHash = hashedPassword

	res, err = a.repository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *authService) Login(ctx context.Context, w model.User) error {
	panic("unimplemented")
}
