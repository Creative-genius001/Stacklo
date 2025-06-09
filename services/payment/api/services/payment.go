package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Creative-genius001/Stacklo/services/payment/config"
	"github.com/Creative-genius001/Stacklo/services/payment/types"
	"github.com/Creative-genius001/go-logger"
)

func GetBankList() (*types.Banks, error) {
	PAYSTACK_BASE_URL := config.Cfg.PaystackBaseUrl
	PAYSTACK_API_KEY := config.Cfg.PaystackTestKey

	// Create client with timeout and retry
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	//create request
	req, err := http.NewRequest("GET", PAYSTACK_BASE_URL+"/bank?currency=NGN", nil)
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

	var banks types.Banks
	if err := json.NewDecoder(resp.Body).Decode(&banks); err != nil {
		logger.Error("failed to decode response: ", err)
		return nil, fmt.Errorf("failed to decode API response")
	}

	return &banks, nil
}

func ResolveAccountNumber(accountNumber string, bankCode string) (*types.AccountResolutionResponse, error) {
	PAYSTACK_BASE_URL := config.Cfg.PaystackBaseUrl
	PAYSTACK_API_KEY := config.Cfg.PaystackTestKey

	// Create client with timeout and retry
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	//create request
	req, err := http.NewRequest("GET", PAYSTACK_BASE_URL+"/bank/resolve?account_number="+accountNumber+"&bank_code="+bankCode, nil)
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

	var acctDetails types.AccountResolutionResponse
	if err := json.NewDecoder(resp.Body).Decode(&acctDetails); err != nil {
		logger.Error("failed to decode response: ", err)
		return nil, fmt.Errorf("failed to decode API response")
	}

	return &acctDetails, nil
}

func CreateTransferRecipient(transferRecipient *types.CreateTransferRecipientRequest) (*types.CreateTransferRecipientResponse, error) {
	PAYSTACK_BASE_URL := config.Cfg.PaystackBaseUrl
	PAYSTACK_API_KEY := config.Cfg.PaystackTestKey

	// Create client with timeout and retry
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	transferRecipientJSON, err := json.Marshal(transferRecipient)
	if err != nil {
		logger.Error("Failed to marshal customer request: ", err)
		return nil, fmt.Errorf("failed to prepare customer data")
	}

	//create request
	req, err := http.NewRequest("POST", PAYSTACK_BASE_URL+"/transferrecipient", bytes.NewBuffer(transferRecipientJSON))
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

	var tRecipient types.CreateTransferRecipientResponse
	if err := json.NewDecoder(resp.Body).Decode(&tRecipient); err != nil {
		logger.Error("failed to decode response: ", err)
		return nil, fmt.Errorf("failed to decode API response")
	}

	return &tRecipient, nil
}

func RequestOTP(payload *types.TransferOtpRequest) (*types.TransferOtpResponse, error) {
	PAYSTACK_BASE_URL := config.Cfg.PaystackBaseUrl
	PAYSTACK_API_KEY := config.Cfg.PaystackTestKey

	// Create client with timeout and retry
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Failed to marshal customer request: ", err)
		return nil, fmt.Errorf("failed to prepare customer data")
	}

	//create request
	req, err := http.NewRequest("POST", PAYSTACK_BASE_URL+"/transfer", bytes.NewBuffer(payloadJSON))
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

	var transfer types.TransferOtpResponse
	if err := json.NewDecoder(resp.Body).Decode(&transfer); err != nil {
		logger.Error("failed to decode response: ", err)
		return nil, fmt.Errorf("failed to decode API response")
	}

	return &transfer, nil
}

func FinalTransfer(payload *types.QueuedTransferRequest) (*types.QueuedTransferResponse, error) {
	PAYSTACK_BASE_URL := config.Cfg.PaystackBaseUrl
	PAYSTACK_API_KEY := config.Cfg.PaystackTestKey

	// Create client with timeout and retry
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Failed to marshal customer request: ", err)
		return nil, fmt.Errorf("failed to prepare customer data")
	}

	//create request
	req, err := http.NewRequest("POST", PAYSTACK_BASE_URL+"/transfer", bytes.NewBuffer(payloadJSON))
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

	var finalTransfer types.QueuedTransferResponse
	if err := json.NewDecoder(resp.Body).Decode(&finalTransfer); err != nil {
		logger.Error("failed to decode response: ", err)
		return nil, fmt.Errorf("failed to decode API response")
	}

	return &finalTransfer, nil
}
