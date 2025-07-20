package service

import (
	"context"
	"errors"
	"time"

	"github.com/Creative-genius001/Stacklo/services/user/config"
	"github.com/Creative-genius001/Stacklo/services/user/email"
	"github.com/Creative-genius001/Stacklo/services/user/utils"
)

var emailService = email.NewEmailClient(config.Cfg.ResendAPI)

type otpService struct {
	repository Repository
}

func (r *otpService) VerifyOTP(ctx context.Context, userID string, otpInput string) error {
	otpRecord, _ := r.repository.GetOTP(ctx, userID)

	if otpRecord.ExpiresAt.Before(time.Now()) {
		return errors.New("OTP expired")
	}
	if otpRecord.Attempts >= 3 {
		return errors.New("Max attempts reached")
	}
	if otpRecord.Verified {
		return errors.New("Already verified")
	}
	if otpRecord.OTP != otpInput {
		r.repository.UpdateOTPCountAttempt(ctx, otpRecord.Attempts+1)
		return errors.New("Invalid OTP")
	}

	return r.repository.UpdateOTPVerificationStatus(ctx, userID, true)
}

func SendOTP(toEmail string) error {
	verificationCode := utils.GenerateOTP(6)
	return emailService.SendVerificationCode(toEmail, verificationCode)
}
