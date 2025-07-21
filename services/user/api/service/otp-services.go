package service

import (
	"context"
	er "errors"
	"time"

	"github.com/Creative-genius001/Stacklo/services/user/email"
	"github.com/Creative-genius001/Stacklo/services/user/redis"
	"github.com/Creative-genius001/Stacklo/services/user/utils"
	errors "github.com/Creative-genius001/Stacklo/services/user/utils/error"
)

type OTPServ interface {
	VerifyOTP(ctx context.Context, userID string, email string, otpInput string) error
	SendOTP(toEmail string) error
}
type otpService struct {
	repository   Repository
	redis        redis.Redis
	emailService email.Resend
}

func NewOTPService(r Repository, rd redis.Redis, e email.Resend) OTPServ {
	return &otpService{r, rd, e}
}

func (r *otpService) VerifyOTP(ctx context.Context, userID string, email string, otpInput string) error {
	otpData, err := r.redis.GetOTPFromRedis(email)
	if err != nil {
		return err
	}
	if otpData == nil || otpData.ExpiresAt.Before(time.Now()) {
		return errors.Wrap(errors.TypeForbidden, "OTP has Expired", er.New("OTP expired"))
	}
	if otpData.Retry >= 3 {
		return errors.Wrap(errors.TypeForbidden, "Max attempts reached", er.New("Max attempts reached"))
	}
	if otpData.OTP != otpInput {
		r.redis.IncrementRetries(email)
		return errors.Wrap(errors.TypeUnauthorized, "Invalid OTP", er.New("Invalid OTP"))
	}

	return r.repository.UpdateVerificationStatus(ctx, userID, true)
}

func (r *otpService) SendOTP(toEmail string) error {
	verificationCode := utils.GenerateOTP(6)
	err := r.redis.SaveOTPToRedis(toEmail, verificationCode)
	if err != nil {
		return errors.Wrap(errors.TypeInternal, "Unable to save OTP to Redis", err)
	}
	return r.emailService.SendVerificationCode(toEmail, verificationCode)
}
