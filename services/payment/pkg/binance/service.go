package binance

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	errors "github.com/Creative-genius001/Stacklo/services/payment/utils/error"
)

func (b *binanceClient) Ping() error {
	_, err := b.BinanceAPIClient(http.MethodGet, "/ping", nil, true)
	if err != nil {
		return err
	}
	return nil
}

func (b *binanceClient) Convert(creq ConvertAssetRequest) ([]*ConvertAssetResponse, error) {

	params := url.Values{}
	params.Set("fromAsset", creq.FromAsset)
	params.Set("toAsset", creq.ToAsset)

	path := "/sapi/v1/convert/exchangeInfo?" + params.Encode()

	resp, err := b.BinanceAPIClient(http.MethodGet, path, nil, true)
	if err != nil {
		return nil, err
	}

	var result []ConvertAssetResponse
	err = json.Unmarshal(resp, &result)
	if err != nil {
		// logger.Logger.Warn("Failed to unmarshal request body", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Failed to unmarshal request body", err)
	}

	res := make([]*ConvertAssetResponse, len(result))
	for i := range result {
		res[i] = &result[i]
	}

	return res, nil
}

func (b *binanceClient) PlaceBuyOrder(ctx context.Context, req BinanceOrderRequest) (*BinanceOrderResponse, error) {
	panic("unimplemented")
}
