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
	"go.uber.org/zap"
)

type Binance interface {
	Order(ctx context.Context, req BinanceOrderRequest) (*BinanceOrderResponse, error)
	Ping() error
	TickerPrice(creq TickerPriceRequest) (*TickerPriceResponse, error)
}

type binanceClient struct {
	client  *http.Client
	logger  *zap.Logger
	baseURL string
	apiKey  string
	secret  string
}

func NewBinanceClient(cfg *config.Config, logger *zap.Logger) Binance {
	return &binanceClient{
		client: &http.Client{
			Timeout: 20 * time.Second,
		},
		logger:  logger,
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
		b.logger.Error("Failed to connect to BINANCE API", zap.Error(err))
		return nil, errors.Wrap(errors.TypeExternal, "Failed to connect to create a new binance request", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if signed || b.apiKey != "" {
		req.Header.Set("X-MBX-APIKEY", b.apiKey)
	}

	b.logger.Debug("BINANCE API CALL", zap.String("path", urlPath), zap.String("method", method))

	resp, err := b.client.Do(req)
	if err != nil {
		b.logger.Error("Failed to get data", zap.Error(err))
		return nil, errors.Wrap(errors.TypeExternal, "Failed to connect to Binance", err)
	}
	defer resp.Body.Close()

	convData, _ := io.ReadAll(resp.Body)
	b.logger.Info("API RESPONSE =>", zap.Any("resp", resp.Body), zap.Any("data", body))

	return convData, nil
}
