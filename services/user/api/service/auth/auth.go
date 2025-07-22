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
	Register(ctx context.Context, user model.User) error
	Login(ctx context.Context, email string, password string) (*model.User, error)
	SignupOTPVerification(ctx context.Context, email string, otp string) error
	AuthOTPVerification(ctx context.Context, email string, otp string) error
}

type authService struct {
	repository   service.Repository
	otp          service.OTPServ
	emailService email.Resend
}

func NewAuthService(r service.Repository, o service.OTPServ, e email.Resend) Auth {
	return &authService{r, o, e}
}

func (a *authService) SignupOTPVerification(ctx context.Context, email string, otp string) error {
	err := a.otp.VerifyOTP(ctx, email, otp)
	if err != nil {
		return err
	}

	err = a.repository.UpdateVerificationStatus(ctx, email, true)
	if err != nil {
		return err
	}
	return nil
}

func (a *authService) AuthOTPVerification(ctx context.Context, email string, otp string) error {
	err := a.otp.VerifyOTP(ctx, email, otp)
	if err != nil {
		return err
	}
	return nil
}

func (a *authService) Register(ctx context.Context, user model.User) error {

	//validate email
	isValid := utils.IsValidEmail(user.Email)
	if !isValid {
		logger.Logger.Warn("invalid email address")
		return errors.Wrap(errors.TypeInvalidInput, "Email address is invalid", er.New("Email address is invalid"))
	}

	//validate phoneNumber
	isValid, formatPhone, err := utils.IsValidPhoneNumber(user.PhoneNumber, "NG")
	if !isValid || err != nil {
		logger.Logger.Warn("invalid phone number")
		return errors.Wrap(errors.TypeInvalidInput, "Phone number is invalid", er.New("Phone number is invalid"))
	}
	user.PhoneNumber = formatPhone

	//validate password
	isValidPassword := utils.IsValidPassword(user.PasswordHash)
	if !isValidPassword {
		logger.Logger.Warn("invalid password")
		return errors.Wrap(errors.TypeInvalidInput, "password must contain uppercase, lowercase and special symbols", er.New("invalid password"))
	}

	resp, err := a.repository.FindByPhoneOrEmail(ctx, user.Email, user.PhoneNumber)
	if err != nil {
		return err
	}
	if resp != nil {
		if resp.Email == user.Email {
			logger.Logger.Warn("user with this email already exists")
			return errors.Wrap(errors.TypeConflict, "email already exists", er.New("email already exists"))
		}
		if resp.PhoneNumber == user.PhoneNumber {
			logger.Logger.Warn("user with this phone number already exists")
			return errors.Wrap(errors.TypeConflict, "phone number already exists", er.New("phone number already exists"))
		}
	}

	hashedPassword, _ := utils.HashPassword(user.PasswordHash)
	user.PasswordHash = hashedPassword

	userResp, err := a.repository.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	a.emailService.SendWelcomeEmail(userResp.Email, fmt.Sprintf("%s %s", userResp.FirstName, userResp.LastName))
	a.otp.SendOTP(userResp.Email)

	return nil
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

	token, _ := utils.CreateToken(user.ID)

	data := model.User{
		ID:          user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Country:     user.Country,
		KycStatus:   user.KycStatus,
		Token:       token,
	}

	return &data, nil
}
