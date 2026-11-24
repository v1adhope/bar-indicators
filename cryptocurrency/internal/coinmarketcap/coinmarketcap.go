package coinmarketcap

import "github.com/shopspring/decimal"

type Data struct {
	Status Status                    `json:"status"`
	Data   map[string]Cryptocurrency `json:"data"`
}

type Status struct {
	Timestamp string `json:"timestamp"`
}

type Cryptocurrency struct {
	Symbol string           `json:"symbol"`
	IsFiat int              `json:"is_fiat"`
	Quote  map[string]Quote `json:"quote"`
}

type Quote struct {
	Price            decimal.Decimal `json:"price"`
	PercentChange24h float64         `json:"percent_change_24h"`
}
