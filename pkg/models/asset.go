package models

type Pair struct {
	FromAsset          string `json:"fromAsset"`
	ToAsset            string `json:"toAsset"`
	FromAssetMinAmount string `json:"fromAssetMinAmount"`
	FromAssetMaxAmount string `json:"fromAssetMaxAmount"`
	ToAssetMinAmount   string `json:"toAssetMinAmount"`
	ToAssetMaxAmount   string `json:"toAssetMaxAmount"`
	FromIsBase         bool   `json:"fromIsBase"`
}

type Triade struct {
	Assets [3]string // Symbols (e.g. "BTC", "ETH", "USDT")
	Pairs  []Pair
}
