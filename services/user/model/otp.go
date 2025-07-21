package model

import "time"

type OTPJSON struct {
	OTP       string    `json:"otp"`
	ExpiresAt time.Time `json:"expires_at"`
	Retry     int       `json:"retry"`
}
