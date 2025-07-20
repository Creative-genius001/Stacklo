package model

import "time"

type OTPJSON struct {
	ID        string    `json:"id"`
	Token     string    `json:"token"`
	OTP       string    `json:"otp"`
	ExpiresAt time.Time `json:"expires_at"`
	Verified  bool      `json:"verified"`
	Attempts  int       `json:"attempts"`
	CreatedAt time.Time `json:"created_at"`
}
