package auth

import (
	"context"
	er "errors"
	"fmt"

	"github.com/Creative-genius001/Stacklo/services/user/api/service"
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
	repository   service.Repository
	otp          service.OTPServ
	emailService email.Resend
}

func NewAuthService(r service.Repository, o service.OTPServ, e email.Resend) Auth {
	return &authService{r, o, e}
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
		return nil, errors.Wrap(errors.TypeInvalidInput, "password must contain uppercase, lowercase and special symbols", er.New("invalid password"))
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
	a.emailService.SendWelcomeEmail(userResp.Email, fmt.Sprintf("%s %s", userResp.FirstName, userResp.LastName))
	a.otp.SendOTP(userResp.Email)

	userData := model.User{
		ID:          userResp.ID,
		Email:       userResp.Email,
		FirstName:   userResp.FirstName,
		LastName:    userResp.LastName,
		Country:     userResp.Country,
		IsVerified:  userResp.IsVerified,
		KycStatus:   userResp.KycStatus,
		PhoneNumber: userResp.PhoneNumber,
	}
	return &userData, nil
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
		return nil, errors.Wrap(errors.TypeConflict, "email or password is incorrect", er.New("email or password is incorrect"))
	}

	return user, nil
}
