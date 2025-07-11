package binance

import (
	"context"

	errors "github.com/Creative-genius001/Stacklo/services/wallet/utils/error"
	"github.com/adshao/go-binance/v2"
)

type Binance interface {
	PlaceOrder(ctx context.Context, req BinanceOrderRequest) (*BinanceOrderResponse, errors.CustomError)
	Ping()
}

type binanceClient struct {
	baseURL       string
	binanceClient *binance.Client
}

func NewBinanceClient(baseURL string, binance *binance.Client) Binance {
	return &binanceClient{
		baseURL:       baseURL,
		binanceClient: binance,
	}
}

func (b *binanceClient) Ping() {

}

func (b *binanceClient) PlaceOrder(ctx context.Context, req BinanceOrderRequest) (*BinanceOrderResponse, errors.CustomError) {
	panic("unimplemented")

}
