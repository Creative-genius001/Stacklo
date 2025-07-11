package binance

import (
	"context"
	"net/http"
)

func (b *binanceClient) Ping() error {
	_, err := b.BinanceAPIClient(http.MethodGet, "ping", nil, true)
	if err != nil {
		return err
	}
	return nil
}

func (b *binanceClient) PlaceBuyOrder(ctx context.Context, req BinanceOrderRequest) (*BinanceOrderResponse, error) {
	panic("unimplemented")
}
