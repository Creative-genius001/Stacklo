package auth

import (
	"context"
	er "errors"

	"github.com/Creative-genius001/Stacklo/services/user/api/service"
	"github.com/Creative-genius001/Stacklo/services/user/config"
	"github.com/Creative-genius001/Stacklo/services/user/email"
	"github.com/Creative-genius001/Stacklo/services/user/model"
	"github.com/Creative-genius001/Stacklo/services/user/utils"
	errors "github.com/Creative-genius001/Stacklo/services/user/utils/error"
	"github.com/Creative-genius001/Stacklo/services/user/utils/logger"
)

type Auth interface {
	Register(ctx context.Context, user model.User) (*model.User, error)
	Login(ctx context.Context, email string, password string) (*model.User, error)
}

type authService struct {
	repository service.Repository
	otp        service.OTPServ
}

func NewAuthService(r service.Repository, o service.OTPServ) Auth {
	return &authService{r, o}
}

var emailService = email.NewEmailClient(config.Cfg.ResendAPI)

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
		return nil, errors.Wrap(errors.TypeInvalidInput, "invalid password", er.New("invalid password"))
	}

	resp, err := a.repository.FindByPhoneOrEmail(ctx, user.Email, user.PhoneNumber)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		if resp.Email == user.Email {
			logger.Logger.Warn("user with this email already exists")
			return nil, errors.Wrap(errors.TypeConflict, "user with this email already exists", er.New("email already exists"))
		}
		if resp.PhoneNumber == user.PhoneNumber {
			logger.Logger.Warn("user with this phone number already exists")
			return nil, errors.Wrap(errors.TypeConflict, "user with this phone number already exists", er.New("phone number already exists"))
		}
	}

	hashedPassword, _ := utils.HashPassword(user.PasswordHash)
	user.PasswordHash = hashedPassword

	userResp, err := a.repository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	a.otp.SendOTP(userResp.Email)
	return userResp, nil
}

func (a *authService) Login(ctx context.Context, email string, password string) (*model.User, error) {

	user, err := a.repository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.Wrap(errors.TypeNotFound, "user not found", er.New("user not found"))
	}

	isValid := utils.CheckPasswordHash(password, user.PasswordHash)
	if !isValid {
		return nil, errors.Wrap(errors.TypeConflict, "password is incorrect", er.New("email or password is incorrect"))
	}

	return user, nil
}
