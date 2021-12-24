package exchangeapi

import "time"

// TickerPrice struct data from exchange source
type TickerPrice struct {
	Ticker   Ticker    `json:"ticker"`
	Time     time.Time `json:"time"`
	Price    string    `json:"price"` // decimal value. example: "0", "10", "12.2", "13.2345122"
	SourceID int       `json:"source_id"`
}

// Ticker is ticker name
type Ticker string

const (
	BTCUSDTicker Ticker = "BTC_USD"
)
