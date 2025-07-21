package model

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	IsVerified   bool      `json:"isVerified"`
	PhoneNumber  string    `json:"phone_number"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Country      string    `json:"country"`
	KycStatus    string    `json:"kyc_status"` //'not_started', 'pending', 'approved', 'rejected', 'resubmit_required'
}
