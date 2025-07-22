package email

import (
	errors "github.com/Creative-genius001/Stacklo/services/user/utils/error"
	"github.com/Creative-genius001/Stacklo/services/user/utils/logger"
	"github.com/resend/resend-go/v2"
	"go.uber.org/zap"
)

type Resend interface {
	SendWelcomeEmail(toEmail string, name string) error
	SendVerificationCode(toEmail string, code string) error
}
type resendClient struct {
	client *resend.Client
}

func NewEmailClient(apiKey string) Resend {
	return &resendClient{
		client: resend.NewClient(apiKey),
	}
}

func (e *resendClient) SendWelcomeEmail(toEmail string, name string) error {
	params := &resend.SendEmailRequest{
		From:    "Stacklo <no-reply@stacklo.com>",
		To:      []string{toEmail},
		Subject: "Welcome to Stacklo",
		Html:    "<strong>Hey " + name + ",</strong><br>Welcome to Stacklo! We're glad you're here.",
	}
	_, err := e.client.Emails.Send(params)
	if err != nil {
		logger.Logger.Error("Failed to send welcome email", zap.Error(err))
		return errors.Wrap(errors.TypeExternal, "Failed to send welcome email", err)
	}
	return nil
}

func (e *resendClient) SendVerificationCode(toEmail string, code string) error {
	params := &resend.SendEmailRequest{
		From:    "Stacklo <no-reply@stacklo.com>",
		To:      []string{toEmail},
		Subject: "Your Stacklo Verification Code",
		Html:    "<p>Your verification code is: <strong>" + code + "</strong></p>",
	}
	_, err := e.client.Emails.Send(params)
	if err != nil {
		logger.Logger.Error("Failed to send otp verification code", zap.Error(err))
		return errors.Wrap(errors.TypeExternal, "Failed to send verification code. Try again", err)
	}
	return nil
}
