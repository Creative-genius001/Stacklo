package paystack

import (
	"bytes"
	"encoding/json"
	er "errors"
	"io"
	"net/http"
	"time"

	errors "github.com/Creative-genius001/Stacklo/services/wallet/utils/error"
	"github.com/Creative-genius001/Stacklo/services/wallet/utils/logger"
	"go.uber.org/zap"
)

type Paystack interface {
	GetBankList() (*Banks, error)
	GetOTP(payload *TransferOtpRequest) (*TransferOtpResponse, error)
	PaystackAPIWrapper(method string, url string, addHeaders map[string]string, data map[string]interface{}) (map[string]interface{}, error)
	ResolveAccountNumber(accountNumber string, bankCode string) (*AccountResolutionResponse, error)
	CreateTransferRecipient(payload *CreateTransferRecipientRequest) (*CreateTransferRecipientResponse, error)
	Transfer(payload FianlTransferRequest) (*FinalTransferResponse, error)
}

type paystackClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewPaystackClient(apiKey, baseURL string) Paystack {
	return &paystackClient{
		apiKey:     apiKey,
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

func (p *paystackClient) GetBankList() (*Banks, error) {

	resp, err := p.PaystackAPIWrapper(http.MethodGet, string(UBankList), nil, nil)
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

	var result Banks
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		logger.Logger.Warn("Failed to unmarshal request body", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Failed to unmarshal request body", err)
	}

	return &result, nil
}

func (p *paystackClient) ResolveAccountNumber(accountNumber string, bankCode string) (*AccountResolutionResponse, error) {
	url := string(UResolveAccNum) + accountNumber + "&bank_code=" + bankCode

	resp, err := p.PaystackAPIWrapper(http.MethodGet, url, nil, nil)
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

	var result AccountResolutionResponse
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		logger.Logger.Warn("Failed to unmarshal request body", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Failed to unmarshal request body", err)
	}

	return &result, nil
}

func (p *paystackClient) CreateTransferRecipient(payload *CreateTransferRecipientRequest) (*CreateTransferRecipientResponse, error) {

	var body map[string]interface{}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &body)

	resp, err := p.PaystackAPIWrapper(http.MethodPost, string(UCreateTrfRecpt), nil, body)
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

	var result CreateTransferRecipientResponse
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		logger.Logger.Warn("Failed to unmarshal request body", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Failed to unmarshal request body", err)
	}

	return &result, nil
}

func (p *paystackClient) GetOTP(payload *TransferOtpRequest) (*TransferOtpResponse, error) {
	var body map[string]interface{}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &body)

	resp, err := p.PaystackAPIWrapper(http.MethodPost, string(UTransfer), nil, body)
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

	var result TransferOtpResponse
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		logger.Logger.Warn("Failed to unmarshal request body", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Failed to unmarshal request body", err)
	}

	return &result, nil
}

func (p *paystackClient) Transfer(payload FianlTransferRequest) (*FinalTransferResponse, error) {
	var body map[string]interface{}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &body)

	resp, err := p.PaystackAPIWrapper(http.MethodPost, string(UFTransfer), nil, body)
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

	var result FinalTransferResponse
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		logger.Logger.Warn("Failed to unmarshal request body", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Failed to unmarshal request body", err)
	}

	return &result, nil
}

func (p *paystackClient) PaystackAPIWrapper(method string, url string, addHeaders map[string]string, data map[string]interface{}) (map[string]interface{}, error) {
	var reqBody io.Reader

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

	PAYSTACK_BASE_URL := p.baseURL
	PAYSTACK_API_KEY := p.apiKey

	req, err := http.NewRequest(method, PAYSTACK_BASE_URL+url, reqBody)
	if err != nil {
		logger.Logger.Error("Failed to connect to PAYSTACK API", zap.Error(err))
		return nil, errors.Wrap(errors.TypeExternal, "Failed to connect to PAYSTACK API", err)
	}

	headers := map[string]string{
		"Authorization": "Bearer " + PAYSTACK_API_KEY,
		"Content-Type":  "application/json",
	}

	for k, v := range addHeaders {
		headers[k] = v
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	var resp *http.Response
	resp, err = p.httpClient.Do(req)
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
		logger.Logger.Error("PAYSTACK API Error Response", zap.Any("method", method), zap.Any("url", url), zap.Int("code", resp.StatusCode), zap.String("error", string(errorBody)))
		return nil, errors.Wrap(errors.TypeExternal, "PAYSTACK API Error", er.New(string(errorBody)))
	}

	var decodedResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&decodedResp); err != nil {
		logger.Logger.Warn("Failed to decode API response", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Failed to decode API response", err)
	}

	return decodedResp, nil
}
