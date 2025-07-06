package types

import "time"

type CreateCustomerRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
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

// Error implements error.
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

type Wallet struct {
	ID                   string    `json:"id,omitempty"`
	Active               bool      `json:"active"`
	UserId               string    `json:"user_id"`
	Currency             string    `json:"currency"`
	Balance              float64   `json:"balance,omitempty"`
	VirtualAccountName   string    `json:"virtual_account_name"`
	VirtualAccountNumber string    `json:"virtual_account_number"`
	VirtualBankName      string    `json:"virtual_bank_name"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}
