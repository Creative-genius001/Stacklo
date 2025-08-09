package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Creative-genius001/Stacklo/services/payment/utils"
	errors "github.com/Creative-genius001/Stacklo/services/payment/utils/error"
	"go.uber.org/zap"
)

func (b *binanceClient) Ping() error {
	_, err := b.BinanceAPIClient(http.MethodGet, "/ping", nil, true)
	if err != nil {
		b.logger.Error("Ping was not successful", zap.Error(err))
		return err
	}
	return nil
}

func (b *binanceClient) TickerPrice(creq TickerPriceRequest) (*TickerPriceResponse, error) {

	params := url.Values{}
	symbol := fmt.Sprintf("%s%s", creq.FromAsset, creq.ToAsset)
	params.Set("symbol", symbol)

	path := "/ticker/price?" + params.Encode()

	resp, err := b.BinanceAPIClient(http.MethodGet, path, nil, true)
	if err != nil {
		return nil, err
	}

	var result TickerPriceResponse
	err = json.Unmarshal(resp, &result)
	if err != nil {
		b.logger.Warn("Failed to unmarshal request body", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Failed to unmarshal request body", err)
	}

	return &result, nil
}

func (b *binanceClient) Order(ctx context.Context, req BinanceOrderRequest) (*BinanceOrderResponse, error) {
	params := url.Values{}
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	params.Set("symbol", req.Symbol)
	params.Set("side", req.Side)
	params.Set("type", "MARKET")
	params.Set("quantity", req.Quantity)
	params.Set("timestamp", timestamp)

	query := params.Encode()
	signature := utils.GenerateHMAC(query)
	fullQuery := query + "&signature=" + signature

	path := "/order?" + fullQuery

	resp, err := b.BinanceAPIClient(http.MethodPost, path, nil, true)
	if err != nil {
		return nil, err
	}

	var result BinanceOrderResponse
	err = json.Unmarshal(resp, &result)
	if err != nil {
		b.logger.Warn("Failed to unmarshal request body", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Failed to unmarshal request body", err)
	}

	return &result, nil
}
