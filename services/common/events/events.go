package events

import (
	"time"

	"github.com/google/uuid"
)

type TransferCompletedEvent struct {
	TransactionID uuid.UUID `json:"transaction_id"` // Unique ID for this transaction
	UserID        uuid.UUID `json:"user_id"`
	WalletID      uuid.UUID `json:"wallet_id"` // The primary wallet affected (e.g., NGN wallet for a withdrawal)
	Amount        float64   `json:"amount"`    // Positive amount
	Currency      string    `json:"currency"`
	EntryType     string    `json:"entry_type"` // "credit" or "debit"
	Status        string    `json:"status"`     // "successful", "failed", "pending" etc.
	Description   string    `json:"description"`
	Timestamp     time.Time `json:"timestamp"`

	// Optional: Details for linking to fiat/crypto transaction details tables
	FiatDetailsID   *uuid.UUID `json:"fiat_details_id,omitempty"`
	CryptoDetailsID *uuid.UUID `json:"crypto_details_id,omitempty"`

	// Additional fields specific to the transfer type
	// For Fiat Transfers (Withdrawals/Deposits)
	GatewayReference *string `json:"gateway_reference,omitempty"` // Paystack ref
	RecipientCode    *string `json:"recipient_code,omitempty"`    // Paystack recipient code
	BankName         *string `json:"bank_name,omitempty"`
	AccountNumber    *string `json:"account_number,omitempty"`
	// ... other fiat specific fields

	// For Crypto Transfers (Buy/Sell/Send/Receive)
	Exchange                  *string  `json:"exchange,omitempty"`          // e.g., "binance"
	ExchangeOrderID           *string  `json:"exchange_order_id,omitempty"` // Binance order ID
	BlockchainTransactionHash *string  `json:"blockchain_transaction_hash,omitempty"`
	Network                   *string  `json:"network,omitempty"` // e.g., "ERC20"
	PriceAtTransaction        *float64 `json:"price_at_transaction,omitempty"`
	QuoteCurrencyAmount       *float64 `json:"quote_currency_amount,omitempty"` // Fiat amount for crypto buy/sell
	// ... other crypto specific fields
}

// You might define other event types here, e.g.:
// type UserRegisteredEvent struct { ... }
// type DepositReceivedEvent struct { ... }
// type WithdrawalInitiatedEvent struct { ... }
