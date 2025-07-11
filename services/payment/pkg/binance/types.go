package binance

type OrderType string

const (
	OrderTypeMarket OrderType = "MARKET"
	OrderTypeLimit  OrderType = "LIMIT"
)

type BinanceOrderRequest struct {
	Symbol        string    // e.g., "BTCUSDT", "ETHNGN"
	Side          string    // "BUY" or "SELL"
	Type          OrderType // "MARKET" or "LIMIT"
	Quantity      float64   // Amount of crypto to buy/sell
	Price         float64   // Optional, for limit orders
	QuoteOrderQty float64   // Optional, for market buys in quote currency (e.g., NGN amount)
}

type BinanceOrderResponse struct {
	OrderID             string  `json:"orderId"`
	Symbol              string  `json:"symbol"`
	Status              string  `json:"status"`              // e.g., "FILLED", "PARTIALLY_FILLED", "NEW"
	ExecutedQty         float64 `json:"executedQty"`         // Amount of crypto filled
	CummulativeQuoteQty float64 `json:"cummulativeQuoteQty"` // Total NGN spent
	ClientOrderID       string  `json:"clientOrderId"`
}
