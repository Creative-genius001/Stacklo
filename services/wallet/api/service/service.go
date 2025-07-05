package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Creative-genius001/Stacklo/services/wallet/config"
	"github.com/Creative-genius001/Stacklo/services/wallet/types"
	"github.com/Creative-genius001/go-logger"
)

type Service interface {
	GetWallet(ctx context.Context, id string) (*types.Wallet, error)
}

type walletService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &walletService{r}
}

func CreateCustomer(customerReq types.CreateCustomerRequest) (*types.CreateCustomerResponse, error) {
	PAYSTACK_BASE_URL := config.Cfg.PaystackBaseUrl
	PAYSTACK_API_KEY := config.Cfg.PaystackTestKey

	customerReqJSON, err := json.Marshal(customerReq)
	if err != nil {
		logger.Error("Failed to marshal customer request: ", err)
		return nil, fmt.Errorf("failed to prepare customer data")
	}

	// Create client with timeout and retry
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	//create request
	req, err := http.NewRequest("POST", PAYSTACK_BASE_URL+"/customer", bytes.NewBuffer(customerReqJSON))
	if err != nil {
		logger.Error("Request creation failed: ", err)
		return nil, fmt.Errorf("failed to create request")
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+PAYSTACK_API_KEY)

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		logger.Error("Failed to make request: ", err)
		return nil, fmt.Errorf("failed to send request")
	}

	defer resp.Body.Close()

	// Handle response
	if resp.StatusCode >= 400 {
		errorBody, _ := io.ReadAll(resp.Body)
		logger.Error("API error" + fmt.Sprint(resp.StatusCode) + ":" + string(errorBody))
		return nil, fmt.Errorf("API error: %s", string(errorBody))
	}

	var customer types.CreateCustomerResponse
	if err := json.NewDecoder(resp.Body).Decode(&customer); err != nil {
		logger.Error("failed to decode response: ", err)
		return nil, fmt.Errorf("failed to decode API response")
	}

	return &customer, nil
}

func CreateDVAWallet(createWalletReq *types.CreateDVAWalletRequest) (*types.CreateDVAWalletResponse, error) {
	var wallet types.CreateDVAWalletResponse

	PAYSTACK_BASE_URL := config.Cfg.PaystackBaseUrl
	PAYSTACK_API_KEY := config.Cfg.PaystackTestKey

	createWalletReqJSON, err := json.Marshal(createWalletReq)
	if err != nil {
		logger.Error("Failed to marshal customer request: ", err)
		return nil, fmt.Errorf("failed to prepare customer data")
	}

	// Create client with timeout and retry
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	//create request
	req, err := http.NewRequest("POST", PAYSTACK_BASE_URL+"/dedicated_account", bytes.NewBuffer(createWalletReqJSON))
	if err != nil {
		logger.Error("Request creation failed: ", err)
		return nil, fmt.Errorf("failed to create request")
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+PAYSTACK_API_KEY)

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		logger.Error("Failed to make request: ", err)
		return nil, fmt.Errorf("failed to send request")
	}

	defer resp.Body.Close()

	// Handle response
	if resp.StatusCode >= 400 {
		errorBody, _ := io.ReadAll(resp.Body)
		logger.Error("API error" + fmt.Sprint(resp.StatusCode) + ":" + string(errorBody))
		return nil, fmt.Errorf("API error: %s", string(errorBody))
	}

	if err := json.NewDecoder(resp.Body).Decode(&wallet); err != nil {
		logger.Error("failed to decode response: ", err)
		return nil, fmt.Errorf("failed to decode API response")
	}

	//save to wallet db

	return &wallet, nil
}

func (w *walletService) GetWallet(ctx context.Context, id string) (*types.Wallet, error) {
	wallet, err := w.GetWallet(ctx, id)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return wallet, nil
}
