package paystack

import (
	"bytes"
	er "errors"
	"fmt"
	"io"
	"net/http"
	"time"

	errors "github.com/Creative-genius001/Stacklo/services/payment/utils/error"
	"github.com/Creative-genius001/Stacklo/services/wallet/utils/logger"
	"go.uber.org/zap"
)

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

func (p *paystackClient) PaystackAPIWrapper(method string, url string, addHeaders map[string]string, data []byte) ([]byte, error) {
	var reqBody io.Reader

	if data != nil {
		reqBody = bytes.NewBuffer(data)
	} else {
		reqBody = nil
	}

	PAYSTACK_API_KEY := p.apiKey
	urlPath := fmt.Sprintf("%s%s", p.baseURL, url)

	req, err := http.NewRequest(method, urlPath, reqBody)
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

	logger.Logger.Debug("API CALL DEBUG", zap.String("path", req.URL.Path), zap.String("method", req.Method))
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

	return io.ReadAll(resp.Body)
}
