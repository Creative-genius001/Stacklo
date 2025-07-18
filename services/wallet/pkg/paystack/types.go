package paystack

import "time"

type URL string

type Paystack interface {
	GetBankList() (*Banks, error)
	GetOTP(payload *TransferOtpRequest) (*TransferOtpResponse, error)
	PaystackAPIWrapper(method string, url string, addHeaders map[string]string, data []byte) ([]byte, error)
	ResolveAccountNumber(accountNumber string, bankCode string) (*AccountResolutionResponse, error)
	CreateTransferRecipient(payload *CreateTransferRecipientRequest) (*CreateTransferRecipientResponse, error)
	Transfer(payload FianlTransferRequest) (*FinalTransferResponse, error)
	CreateCustomer(payload CreateCustomerRequest) (*CreateCustomerResponse, error)
	CreateDVAWallet(createWalletReq *CreateDVAWalletRequest) (*CreateDVAWalletResponse, error)
}

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

type CreateCustomerRequest struct {
	Email     string `json:"email" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
}

type CreateCustomerResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Email           string      `json:"email"`
		Integration     int64       `json:"integration"`
		Domain          string      `json:"domain"`
		CustomerCode    string      `json:"customer_code"`
		ID              int64       `json:"id"`
		Identified      bool        `json:"identified"`
		Identifications interface{} `json:"identifications"`
		CreatedAt       time.Time   `json:"createdAt"`
		UpdatedAt       time.Time   `json:"updatedAt"`
	} `json:"data"`
}

type CreateDVAWalletRequest struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Phone         string `json:"phone"`
	PreferredBank string `json:"preferred_bank"`
	CustomerCode  int64  `json:"customer"`
}

type CreateDVAWalletResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Bank struct {
			Name string `json:"name"`
			ID   int64  `json:"id"`
			Slug string `json:"slug"`
		} `json:"bank"`
		AccountName   string      `json:"account_name"`
		AccountNumber string      `json:"account_number"`
		Assigned      bool        `json:"assigned"`
		Currency      string      `json:"currency"`
		Metadata      interface{} `json:"metadata"`
		Active        bool        `json:"active"`
		ID            int64       `json:"id"`
		CreatedAt     time.Time   `json:"created_at"`
		UpdatedAt     time.Time   `json:"updated_at"`
		Assignment    struct {
			Integration  int64     `json:"integration"`
			AssigneeID   int64     `json:"assignee_id"`
			AssigneeType string    `json:"assignee_type"`
			Expired      bool      `json:"expired"`
			AccountType  string    `json:"account_type"`
			AssignedAt   time.Time `json:"assigned_at"`
		} `json:"assignment"`
		Customer struct {
			ID           int64  `json:"id"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			Email        string `json:"email"`
			CustomerCode string `json:"customer_code"`
			Phone        string `json:"phone"`
			RiskAction   string `json:"risk_action"`
		} `json:"customer"`
	} `json:"data"`
}

func (c CreateCustomerResponse) Error() string {
	panic("unimplemented")
}

type GetWalletResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Customer struct {
			ID           int64  `json:"id"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			Email        string `json:"email"`
			CustomerCode string `json:"customer_code"`
			Phone        string `json:"phone"`
			Metadata     struct {
				CallingCode string `json:"calling_code"`
			} `json:"metadata"`
			RiskAction               string      `json:"risk_action"`
			InternationalFormatPhone interface{} `json:"international_format_phone"`
		} `json:"customer"`
		Bank struct {
			Name string `json:"name"`
			ID   int    `json:"id"`
			Slug string `json:"slug"`
		} `json:"bank"`
		ID            int    `json:"id"`
		AccountName   string `json:"account_name"`
		AccountNumber string `json:"account_number"`
		CreatedAt     string `json:"created_at"`
		UpdatedAt     string `json:"updated_at"`
		Currency      string `json:"currency"`
		SplitConfig   string `json:"split_config"`
		Active        bool   `json:"active"`
		Assigned      bool   `json:"assigned"`
	} `json:"data"`
}
