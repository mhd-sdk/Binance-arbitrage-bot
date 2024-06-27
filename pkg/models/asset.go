package models

type Triade struct {
	Assets  []string // Symbols (e.g. "BTC", "ETH", "USDT")
	Symbols []Symbol
}

// this function return two possible orders for the triade (both starting from the first asset)
func (t Triade) OrderAssets(startingAsset string) ([]string, []string) {
	order1 := []string{}
	order2 := []string{}
	startingIdx := 0
	for c, asset := range t.Assets {
		if asset == startingAsset {
			startingIdx = c
		}
	}
	switch startingIdx {
	case 0:
		order1 = t.Assets
		order2 = []string{t.Assets[0], t.Assets[2], t.Assets[1]}

	case 1:
		order1 = []string{t.Assets[1], t.Assets[0], t.Assets[2]}
		order2 = []string{t.Assets[1], t.Assets[2], t.Assets[0]}
	case 2:
		order1 = []string{t.Assets[2], t.Assets[1], t.Assets[0]}
		order2 = []string{t.Assets[2], t.Assets[0], t.Assets[1]}
	}
	return order1, order2
}

type Trade struct {
	EventType string `json:"e"`
	EventTime int    `json:"E"`
	Symbol    string `json:"s"`
	TradeID   int    `json:"t"`
	Price     string `json:"p"`
	Quantity  string `json:"q"`
	TradeTime int    `json:"T"`
	Maker     bool   `json:"m"`
	BestMatch bool   `json:"M"`
}
