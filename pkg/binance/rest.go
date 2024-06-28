package binance

import (
	"encoding/json"
	"io"
	"math/big"
	"net/http"

	binance_connector "github.com/binance/binance-connector-go"
)

func GetSymbolPrices() (map[string]*big.Float, error) {
	url := "https://api.binance.com/api/v3/ticker/price"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var symbolPrices []*binance_connector.TickerPriceResponse

	err = json.Unmarshal(body, &symbolPrices)
	if err != nil {
		return nil, err
	}

	prices := map[string]*big.Float{}

	for _, symbolPrice := range symbolPrices {
		// price, err := strconv.ParseFloat(symbolPrice.Price, 64)
		price, _, err := big.ParseFloat(symbolPrice.Price, 10, 64, big.ToNearestEven)
		if err != nil {
			return nil, err
		}
		prices[symbolPrice.Symbol] = price
	}

	return prices, nil
}
