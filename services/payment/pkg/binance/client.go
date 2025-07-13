package binance

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"net/http"
	"time"

	"github.com/Creative-genius001/Stacklo/services/payment/config"
	errors "github.com/Creative-genius001/Stacklo/services/payment/utils/error"
	// "github.com/Creative-genius001/Stacklo/services/payment/utils/logger"
	// "go.uber.org/zap"
)

type Binance interface {
	PlaceBuyOrder(ctx context.Context, req BinanceOrderRequest) (*BinanceOrderResponse, error)
	Ping() error
	Convert(creq ConvertAssetRequest) ([]*ConvertAssetResponse, error)
}

type binanceClient struct {
	client  *http.Client
	baseURL string
	apiKey  string
	secret  string
}

func NewBinanceClient(cfg *config.Config) Binance {
	return &binanceClient{
		client: &http.Client{
			Timeout: 20 * time.Second,
		},
		baseURL: cfg.BinanceBaseUrl,
		apiKey:  cfg.BinanceAPIKey,
		secret:  cfg.BinanceSecretKey,
	}
}

// type binanceClient struct {
// 	baseURL       string
// 	binanceClient *binance.Client
// }

// func NewBinanceClient(baseURL string, binance *binance.Client) Binance {
// 	return &binanceClient{
// 		baseURL:       baseURL,
// 		binanceClient: binance,
// 	}
// }

func (b *binanceClient) BinanceAPIClient(method, url string, body any, signed bool) ([]byte, error) {
	urlPath := fmt.Sprintf("%s%s", b.baseURL, url)

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(b)
	} else {
		bodyReader = nil
	}
	req, err := http.NewRequest(method, urlPath, bodyReader)
	if err != nil {
		// logger.Logger.Error("Failed to connect to BINANCE API", zap.Error(err))
		return nil, errors.Wrap(errors.TypeExternal, "Failed to connect to BINANCE API", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if signed || b.apiKey != "" {
		req.Header.Set("X-MBX-APIKEY", b.apiKey)
	}

	// logger.Logger.Debug("API CALL DEBUG", zap.String("path", req.URL.Path), zap.String("method", req.Method))

	resp, err := b.client.Do(req)
	if err != nil {
		// logger.Logger.Error("Failed to get data", zap.Error(err))
		return nil, errors.Wrap(errors.TypeExternal, "Failed to get BINANCE API data", err)
	}
	defer resp.Body.Close()

	// logger.Logger.Info("Binance response body", zap.ByteString("body", respBody))

	return io.ReadAll(resp.Body)
}
