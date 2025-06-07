package services

import (
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
