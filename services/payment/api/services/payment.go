package services

import (
	"bytes"

	"encoding/json"
	er "errors"
	"io"
	"net/http"
	"time"

	"github.com/Creative-genius001/Stacklo/services/payment/config"
	"github.com/Creative-genius001/Stacklo/services/payment/types"
	errors "github.com/Creative-genius001/Stacklo/services/wallet/utils/error"
	"github.com/Creative-genius001/Stacklo/services/wallet/utils/logger"
	"go.uber.org/zap"
)

func PaystackAPIWrapper(method string, url string, addHeaders map[string]string, data map[string]interface{}) (map[string]interface{}, error) {
	PAYSTACK_BASE_URL := config.Cfg.PaystackBaseUrl
	PAYSTACK_API_KEY := config.Cfg.PaystackTestKey
	var reqBody io.Reader

	// Create client with timeout and retry
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	if data != nil {
		json, err := json.Marshal(data)
		if err != nil {
			logger.Logger.Warn("Failed to marshal request body", zap.Error(err))
			return nil, errors.Wrap(errors.TypeInternal, "Failed to marshal request body", err)
		}
		reqBody = bytes.NewBuffer(json)
	} else {
		reqBody = nil
	}

	//create request
	req, err := http.NewRequest(method, PAYSTACK_BASE_URL+url, reqBody)
	if err != nil {
		logger.Logger.Error("Failed to connect to PAYSTACK API", zap.Error(err))
		return nil, errors.Wrap(errors.TypeExternal, "Failed to connect to PAYSTACK API", err)
	}

	//set headers and additional headers by mappping each header to their key value pairs and looping over them
	headers := map[string]string{
		"Authorization": "Bearer " + PAYSTACK_API_KEY,
	}

	for k, v := range addHeaders {
		headers[k] = v
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		logger.Logger.Error("Failed to return PAYSTACK API call response", zap.Any("method", method), zap.Any("url", url), zap.Error(err))
		return nil, errors.Wrap(errors.TypeExternal, "Failed to return PAYSTACK API call response", err)
	}
	defer resp.Body.Close()

	// Handle response
	if resp.StatusCode >= 400 {
		errorBody, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Logger.Warn("Failed to read response body", zap.Error(err))
			return nil, errors.Wrap(errors.TypeExternal, "Failed to read response body", er.New("Failed to read response body"))
		}
		logger.Logger.Error("PAYSTACK API Error Response", zap.Any("method", method), zap.Any("url", url), zap.String("error", string(errorBody)))
		return nil, errors.Wrap(errors.TypeExternal, "PAYSTACK API Error", er.New(string(errorBody)))
	}

	var decodedResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&decodedResp); err != nil {
		logger.Logger.Warn("Failed to decode API response", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Failed to decode API response", err)
	}

	return decodedResp, nil
}

func GetBankList() (*types.Banks, error) {

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	resp, err := PaystackAPIWrapper("GET", string(types.UBankList), headers, nil)
	if err != nil {
		var appErr *errors.CustomError
		if !er.As(err, &appErr) {
			return nil, errors.Wrap(errors.TypeExternal, "External API error", err)
		}

		logger.Logger.Error("It aserted at the service level", zap.String("messg", appErr.Message), zap.String("type", string(appErr.Type)), zap.Error(appErr.Err))
		return nil, appErr
	}

	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		logger.Logger.Warn("Failed to marshal request body", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Failed to marshal request body", err)
	}

	var bankList types.Banks
	err = json.Unmarshal(jsonBytes, &bankList)
	if err != nil {
		logger.Logger.Warn("Failed to unmarshal request body", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Failed to unmarshal request body", err)
	}

	return &bankList, nil
}

func ResolveAccountNumber(accountNumber string, bankCode string) (*types.AccountResolutionResponse, error) {

	url := string(types.UResolveAccNum) + accountNumber + "&bank_code=" + bankCode

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	resp, err := PaystackAPIWrapper("GET", url, headers, nil)
	if err != nil {
		var appErr *errors.CustomError
		if !er.As(err, &appErr) {
			return nil, errors.Wrap(errors.TypeExternal, "External API error", err)
		}

		logger.Logger.Error("It aserted at the service level", zap.String("messg", appErr.Message), zap.String("type", string(appErr.Type)), zap.Error(appErr.Err))
		return nil, appErr
	}

	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		logger.Logger.Warn("Failed to marshal request body", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Failed to marshal request body", err)
	}

	var acctDetails types.AccountResolutionResponse
	err = json.Unmarshal(jsonBytes, &acctDetails)
	if err != nil {
		logger.Logger.Warn("Failed to unmarshal request body", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Failed to unmarshal request body", err)
	}

	return &acctDetails, nil
}

// func CreateTransferRecipient(transferRecipient *types.CreateTransferRecipientRequest) (*types.CreateTransferRecipientResponse, error) {
// 	PAYSTACK_BASE_URL := config.Cfg.PaystackBaseUrl
// 	PAYSTACK_API_KEY := config.Cfg.PaystackTestKey

// 	// Create client with timeout and retry
// 	client := &http.Client{
// 		Timeout: 15 * time.Second,
// 	}

// 	transferRecipientJSON, err := json.Marshal(transferRecipient)
// 	if err != nil {
// 		logger.Error("Failed to marshal customer request: ", err)
// 		return nil, fmt.Errorf("failed to prepare customer data")
// 	}

// 	//create request
// 	req, err := http.NewRequest("POST", PAYSTACK_BASE_URL+"/transferrecipient", bytes.NewBuffer(transferRecipientJSON))
// 	if err != nil {
// 		logger.Error("Request creation failed: ", err)
// 		return nil, fmt.Errorf("failed to create request")
// 	}

// 	// Set headers
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", "Bearer "+PAYSTACK_API_KEY)

// 	var resp *http.Response
// 	resp, err = client.Do(req)
// 	if err != nil {
// 		logger.Error("Failed to make request: ", err)
// 		return nil, fmt.Errorf("failed to send request")
// 	}

// 	defer resp.Body.Close()

// 	// Handle response
// 	if resp.StatusCode >= 400 {
// 		errorBody, _ := io.ReadAll(resp.Body)
// 		logger.Error("API error" + fmt.Sprint(resp.StatusCode) + ":" + string(errorBody))
// 		return nil, fmt.Errorf("API error: %s", string(errorBody))
// 	}

// 	var tRecipient types.CreateTransferRecipientResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&tRecipient); err != nil {
// 		logger.Error("failed to decode response: ", err)
// 		return nil, fmt.Errorf("failed to decode API response")
// 	}

// 	return &tRecipient, nil
// }

// func RequestOTP(payload *types.TransferOtpRequest) (*types.TransferOtpResponse, error) {
// 	PAYSTACK_BASE_URL := config.Cfg.PaystackBaseUrl
// 	PAYSTACK_API_KEY := config.Cfg.PaystackTestKey

// 	// Create client with timeout and retry
// 	client := &http.Client{
// 		Timeout: 15 * time.Second,
// 	}

// 	payloadJSON, err := json.Marshal(payload)
// 	if err != nil {
// 		logger.Error("Failed to marshal customer request: ", err)
// 		return nil, fmt.Errorf("failed to prepare customer data")
// 	}

// 	//create request
// 	req, err := http.NewRequest("POST", PAYSTACK_BASE_URL+"/transfer", bytes.NewBuffer(payloadJSON))
// 	if err != nil {
// 		logger.Error("Request creation failed: ", err)
// 		return nil, fmt.Errorf("failed to create request")
// 	}

// 	// Set headers
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", "Bearer "+PAYSTACK_API_KEY)

// 	var resp *http.Response
// 	resp, err = client.Do(req)
// 	if err != nil {
// 		logger.Error("Failed to make request: ", err)
// 		return nil, fmt.Errorf("failed to send request")
// 	}

// 	defer resp.Body.Close()

// 	// Handle response
// 	if resp.StatusCode >= 400 {
// 		errorBody, _ := io.ReadAll(resp.Body)
// 		logger.Error("API error" + fmt.Sprint(resp.StatusCode) + ":" + string(errorBody))
// 		return nil, fmt.Errorf("API error: %s", string(errorBody))
// 	}

// 	var transfer types.TransferOtpResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&transfer); err != nil {
// 		logger.Error("failed to decode response: ", err)
// 		return nil, fmt.Errorf("failed to decode API response")
// 	}

// 	return &transfer, nil
// }

// func FinalTransfer(payload *types.QueuedTransferRequest) (*types.QueuedTransferResponse, error) {
// 	PAYSTACK_BASE_URL := config.Cfg.PaystackBaseUrl
// 	PAYSTACK_API_KEY := config.Cfg.PaystackTestKey

// 	// Create client with timeout and retry
// 	client := &http.Client{
// 		Timeout: 15 * time.Second,
// 	}

// 	payloadJSON, err := json.Marshal(payload)
// 	if err != nil {
// 		logger.Error("Failed to marshal customer request: ", err)
// 		return nil, fmt.Errorf("failed to prepare customer data")
// 	}

// 	//create request
// 	req, err := http.NewRequest("POST", PAYSTACK_BASE_URL+"/transfer", bytes.NewBuffer(payloadJSON))
// 	if err != nil {
// 		logger.Error("Request creation failed: ", err)
// 		return nil, fmt.Errorf("failed to create request")
// 	}

// 	// Set headers
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", "Bearer "+PAYSTACK_API_KEY)

// 	var resp *http.Response
// 	resp, err = client.Do(req)
// 	if err != nil {
// 		logger.Error("Failed to make request: ", err)
// 		return nil, fmt.Errorf("failed to send request")
// 	}

// 	defer resp.Body.Close()

// 	// Handle response
// 	if resp.StatusCode >= 400 {
// 		errorBody, _ := io.ReadAll(resp.Body)
// 		logger.Error("API error" + fmt.Sprint(resp.StatusCode) + ":" + string(errorBody))
// 		return nil, fmt.Errorf("API error: %s", string(errorBody))
// 	}

// 	var finalTransfer types.QueuedTransferResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&finalTransfer); err != nil {
// 		logger.Error("failed to decode response: ", err)
// 		return nil, fmt.Errorf("failed to decode API response")
// 	}

// 	return &finalTransfer, nil
// }
