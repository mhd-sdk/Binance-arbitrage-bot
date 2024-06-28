package models

type Triade struct {
	Assets  []string // (e.g. "BTC", "ETH", "USDT")
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

// order symbols based on the starting symbol
func (t Triade) OrderSymbols(startingAsset string) ([]Symbol, []Symbol) {
	order1Symbols := []Symbol{}
	order2Symbols := []Symbol{}

	order1Assets, order2Assets := t.OrderAssets(startingAsset)

	for c, _ := range order1Assets {
		from := order1Assets[c]
		to := order1Assets[(c+1)%3]
		order1Symbols = append(order1Symbols, FindSymbol(from, to, t.Symbols))
	}

	for c, _ := range order2Assets {
		from := order2Assets[c]
		to := order2Assets[(c+1)%3]
		order2Symbols = append(order2Symbols, FindSymbol(from, to, t.Symbols))
	}

	return order1Symbols, order2Symbols

}

func FindSymbol(asset1 string, asset2 string, symbols []Symbol) Symbol {
	for _, symbol := range symbols {
		if symbol.Symbol == asset1+asset2 {
			return symbol
		} else if symbol.Symbol == asset2+asset1 {
			return symbol
		}

	}
	return Symbol{}
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
