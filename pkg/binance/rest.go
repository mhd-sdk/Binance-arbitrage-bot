package binance

import (
	"DeltA/pkg/models"
	"encoding/json"
	"io"
	"net/http"
)

type ExchangeInfo []models.Pair

func GetExchangeInfo() (ExchangeInfo, error) {
	url := "https://api.binance.com/sapi/v1/convert/exchangeInfo"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	exchangeInfo := ExchangeInfo{}

	err = json.Unmarshal(body, &exchangeInfo)
	if err != nil {
		return nil, err
	}

	return exchangeInfo, nil
}
