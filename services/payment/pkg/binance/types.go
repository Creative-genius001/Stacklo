package binance

type OrderType string

const (
	OrderTypeMarket OrderType = "MARKET"
	OrderTypeLimit  OrderType = "LIMIT"
)

type BinanceOrderRequest struct {
	Symbol   string `json:"symbol"`   // e.g., "BTCUSDT", "ETHNGN"
	Side     string `json:"side"`     // "BUY" or "SELL"
	Quantity string `json:"quantity"` // Amount of crypto to buy/sell
}

type BinanceOrderResponse struct {
	OrderID             string  `json:"orderId"`
	Symbol              string  `json:"symbol"`
	Status              string  `json:"status"`              // e.g., "FILLED", "PARTIALLY_FILLED", "NEW"
	ExecutedQty         float64 `json:"executedQty"`         // Amount of crypto filled
	CummulativeQuoteQty float64 `json:"cummulativeQuoteQty"` // Total NGN spent
	ClientOrderID       string  `json:"clientOrderId"`
}

type TickerPriceRequest struct {
	FromAsset string `json:"fromAsset"`
	ToAsset   string `json:"toAsset"`
}

type TickerPriceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}
