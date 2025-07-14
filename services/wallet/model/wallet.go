package model

import "time"

type Wallet struct {
	ID                   string    `json:"id,omitempty"`
	Active               bool      `json:"active"`
	UserId               string    `json:"user_id"`
	Currency             string    `json:"currency"`
	WalletType           string    `json:"wallet_type"`
	Balance              float64   `json:"balance,omitempty"`
	VirtualAccountName   *string   `json:"virtual_account_name,omitempty"`
	VirtualAccountNumber *string   `json:"virtual_account_number,omitempty"`
	VirtualBankName      *string   `json:"virtual_bank_name,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}
