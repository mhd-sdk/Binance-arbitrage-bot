package models

type ExchangeInfo struct {
	Timezone string `json:"timezone"`
	Symbols  []Symbol
}

type Symbol struct {
	Symbol               string `json:"symbol"`
	Status               string `json:"status"`
	BaseAsset            string `json:"baseAsset"`
	QuoteAsset           string `json:"quoteAsset"`
	IsSpotTradingAllowed bool   `json:"isSpotTradingAllowed"`
}
