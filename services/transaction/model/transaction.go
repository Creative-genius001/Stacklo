package model

import "time"

type Transaction struct {
	ID              string             `json:"id"`
	UserId          string             `json:"user_id"`
	WalletId        string             `json:"wallet_id"`
	Currency        string             `json:"currency"` //NGN, BTC,ETH,USDT
	Status          float64            `json:"status"`   //PENDING,SUCCESS,FAILED,REVERSED,PROCESSING
	Amount          float64            `json:"amount"`
	Reason          float64            `json:"reason"`
	TransactionType string             `json:"transaction_type"` //CRYPTO,FIAT
	EntryType       string             `json:"entry_type"`       //CREDIT,DEBIT
	FiatDetails     *FiatTransaction   `json:"fiat_details,omitempty"`
	CryptoDetails   *CryptoTransaction `json:"crypto_details,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}

type FiatTransaction struct {
	ID                string  `json:"id"`
	ReferenceID       string  `json:"reference_id"`
	TransactionNumber string  `json:"transaction_number"`
	BankName          string  `json:"bank_name"`
	AccountName       string  `json:"account_name"`
	AccountNumber     string  `json:"account_number"`
	Fee               float64 `json:"fee"`
	NetAmount         float64 `json:"net_amount"`
}

type CryptoTransaction struct {
	ID                  string  `json:"id"`
	ExchangeOrderID     string  `json:"exchange_order_id"`
	Network             string  `json:"network"`
	NetworkFee          float64 `json:"network_fee"`
	PriceAtTransaction  float64 `json:"price_at_transaction"`
	QuoteCurrencyAmount float64 `json:"quote_currency_amount"`
}
