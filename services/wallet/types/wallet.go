package types

import "time"

type CreateCustomerRequest struct {
	Email         string `json:"email"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Phone         string `json:"phone"`
	PreferredBank string `json:"preferred_bank"`
	Country       string `json:"country"`
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
	Email         string `json:"email"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Phone         string `json:"phone"`
	PreferredBank string `json:"preferred_bank"`
	Country       string `json:"country"`
	Customer      int64  `json:"customer"`
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
