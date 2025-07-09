package types

import "context"

type Services interface {
	GetBankList() (*Banks, error)
	PaystackAPIWrapper(ctx context.Context, url string, headers map[string]interface{}, data map[string]interface{}) (map[string]interface{}, error)
	GetOTP(ctx context.Context, trf TransferOtpRequest) (*TransferOtpResponse, error)
}

type URL string

const (
	UBankList       URL = "/bank?currency=NGN"
	UResolveAccNum  URL = "/bank/resolve?account_number="
	UCreateTrfRecpt URL = "/transferrecipient"
	UTransfer       URL = "/transfer"
	UFTransfer      URL = "/transfer/finalize_transfer"
)

type StartTransferData struct {
	Name          string  `json:"name"`
	Type          string  `json:"type,omitempty"`
	Currency      string  `json:"currency,omitempty"`
	AccountNumber string  `json:"account_number"`
	BankCode      string  `json:"bank_code"`
	Reason        string  `json:"reason,omitempty"`
	Amount        float64 `json:"amount"`
}
type CreateTransferRecipientRequest struct {
	Type          string `json:"type"`
	Name          string `json:"name"`
	AccountNumber string `json:"account_number"`
	BankCode      string `json:"bank_code"`
	Currency      string `json:"currency"`
}

type CreateTransferRecipientResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Active        bool   `json:"active"`
		CreatedAt     string `json:"createdAt"`
		Currency      string `json:"currency"`
		Domain        string `json:"domain"`
		ID            int64  `json:"id"`
		Integration   int64  `json:"integration"`
		Name          string `json:"name"`
		RecipientCode string `json:"recipient_code"`
		Type          string `json:"type"`
		UpdatedAt     string `json:"updatedAt"`
		IsDeleted     bool   `json:"is_deleted"`
		Details       struct {
			AuthorizationCode *string `json:"authorization_code"`
			AccountNumber     string  `json:"account_number"`
			AccountName       string  `json:"account_name"`
			BankCode          string  `json:"bank_code"`
			BankName          string  `json:"bank_name"`
		} `json:"details"`
	} `json:"data"`
}

type TransferOtpRequest struct {
	Source    string  `json:"source"`
	Reason    string  `json:"reason"`
	Amount    float64 `json:"amount"`
	Recipeint string  `json:"recipient"`
	Reference string  `json:"reference"`
}

type TransferOtpResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Integration  int64  `json:"integration"`
		Domain       string `json:"domain"`
		Amount       int    `json:"amount"`
		Currency     string `json:"currency"`
		Source       string `json:"source"`
		Reason       string `json:"reason"`
		Recipient    int64  `json:"recipient"`
		Status       string `json:"status"`
		TransferCode string `json:"transfer_code"`
		ID           int64  `json:"id"`
		CreatedAt    string `json:"createdAt"`
		UpdatedAt    string `json:"updatedAt"`
	} `json:"data"`
}

type FianlTransferRequest struct {
	TransferCode string `json:"transfer_code"`
	Otp          string `json:"otp"`
}

type FinalTransferResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Domain        string      `json:"domain"`
		Amount        int         `json:"amount"`
		Currency      string      `json:"currency"`
		Reference     string      `json:"reference"`
		Source        string      `json:"source"`
		SourceDetails interface{} `json:"source_details"`
		Reason        string      `json:"reason"`
		Status        string      `json:"status"`
		Failures      interface{} `json:"failures"`
		TransferCode  string      `json:"transfer_code"`
		TitanCode     *string     `json:"titan_code"`
		TransferredAt *string     `json:"transferred_at"`
		ID            int64       `json:"id"`
		Integration   int64       `json:"integration"`
		Recipient     int64       `json:"recipient"`
		CreatedAt     string      `json:"createdAt"`
		UpdatedAt     string      `json:"updatedAt"`
	} `json:"data"`
}

type FinalTransferJson struct {
	Reference string `json:"reference"`
	Amount    int    `json:"amount"`
	Currency  string `json:"currency"`
	Recipient int64  `json:"recipient"`
	Reason    string `json:"reason"`
	Status    string `json:"status"`
}

type PaystackError struct {
	Status  string
	Message string
	Meta    []map[string]interface{}
}
type Banks struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    []struct {
		Name        string      `json:"name"`
		Slug        string      `json:"slug"`
		Code        string      `json:"code"`
		Longcode    string      `json:"longcode"`
		Gateway     interface{} `json:"gateway"`
		PayWithBank bool        `json:"pay_with_bank"`
		Active      bool        `json:"active"`
		IsDeleted   bool        `json:"is_deleted,omitempty"`
		Country     string      `json:"country,omitempty"`
		Currency    string      `json:"currency,omitempty"`
		Type        string      `json:"type,omitempty"`
		CreatedAt   string      `json:"createdAt,omitempty"`
		UpdatedAt   string      `json:"updatedAt,omitempty"`
	} `json:"data"`
}

type AccountResolutionResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AccountNumber string `json:"account_number"`
		AccountName   string `json:"account_name"`
		BankID        int    `json:"bank_id"`
	} `json:"data"`
}
