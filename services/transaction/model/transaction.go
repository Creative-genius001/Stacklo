package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID              string             `json:"id"`
	UserId          string             `json:"user_id"`
	WalletId        string             `json:"wallet_id"`
	Currency        string             `json:"currency"` //NGN, BTC,ETH,USDT
	Status          string             `json:"status"`   //PENDING,SUCCESS,FAILED,REVERSED,PROCESSING
	Amount          decimal.Decimal    `json:"amount"`
	Reason          string             `json:"reason"`
	TransactionType string             `json:"transaction_type"` //CRYPTO,FIAT
	EntryType       string             `json:"entry_type"`       //CREDIT,DEBIT
	FiatDetails     *FiatTransaction   `json:"fiat_details,omitempty"`
	CryptoDetails   *CryptoTransaction `json:"crypto_details,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}

type FiatTransaction struct {
	ID                *string          `json:"id"`
	ReferenceID       *string          `json:"reference_id"`
	TransactionNumber *string          `json:"transaction_number"`
	BankName          *string          `json:"bank_name"`
	AccountName       *string          `json:"account_name"`
	AccountNumber     *string          `json:"account_number"`
	Fee               *decimal.Decimal `json:"fee"`
	NetAmount         *decimal.Decimal `json:"net_amount"`
}

type CryptoTransaction struct {
	ID                  *string          `json:"id"`
	ExchangeOrderID     *string          `json:"exchange_order_id"`
	Network             *string          `json:"network"`
	NetworkFee          *decimal.Decimal `json:"network_fee"`
	PriceAtTransaction  *decimal.Decimal `json:"price_at_transaction"`
	QuoteCurrencyAmount *decimal.Decimal `json:"quote_currency_amount"`
}
