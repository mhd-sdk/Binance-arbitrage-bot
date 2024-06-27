package binance

import (
	"DeltA/pkg/models"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func GetExchangeInfo() (models.ExchangeInfo, error) {
	url := "https://api.binance.com/api/v3/exchangeInfo"
	resp, err := http.Get(url)
	if err != nil {
		return models.ExchangeInfo{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.ExchangeInfo{}, err
	}

	exchangeInfo := models.ExchangeInfo{}

	err = json.Unmarshal(body, &exchangeInfo)
	if err != nil {
		return models.ExchangeInfo{}, err
	}

	return exchangeInfo, nil
}

func GetSymbolPrices() (map[string]float64, error) {
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

	var symbolPrices []struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}

	err = json.Unmarshal(body, &symbolPrices)
	if err != nil {
		return nil, err
	}

	prices := map[string]float64{}

	for _, symbolPrice := range symbolPrices {
		price, err := strconv.ParseFloat(symbolPrice.Price, 64)
		if err != nil {
			return nil, err
		}
		prices[symbolPrice.Symbol] = price
	}

	return prices, nil
}
