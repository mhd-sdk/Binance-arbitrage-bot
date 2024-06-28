package helpers

import (
	"DeltA/pkg/models"
	"slices"
	"strings"
)

func ListAllAssets(exchangeInfo models.ExchangeInfo) []string {
	availableAssets := []string{}
	for _, pair := range exchangeInfo.Symbols {
		if !slices.Contains(availableAssets, pair.BaseAsset) {
			availableAssets = append(availableAssets, pair.BaseAsset)
		}
		if !slices.Contains(availableAssets, pair.QuoteAsset) {
			availableAssets = append(availableAssets, pair.QuoteAsset)
		}
	}
	return availableAssets
}

func BuildTriades(exchangeInfo models.ExchangeInfo, startingAsset string) []models.Triade {
	allAssets := ListAllAssets(exchangeInfo)
	tempTriades := []models.Triade{}

	// build every possible triade containing the starting asset
	for i := 0; i < len(allAssets); i++ {
		for j := i + 1; j < len(allAssets); j++ {
			for k := j + 1; k < len(allAssets); k++ {
				if allAssets[i] != startingAsset && allAssets[j] != startingAsset && allAssets[k] != startingAsset {
					continue
				}
				tempTriades = append(tempTriades, models.Triade{
					Assets: []string{allAssets[i], allAssets[j], allAssets[k]},
				})
			}
		}
	}

	// check if triades are valid and append symbols
	finalTriades := []models.Triade{}
	for x, triade := range tempTriades {
		tempTriades[x].Symbols = buildTriadeSymbols(exchangeInfo, triade)
		if len(tempTriades[x].Symbols) == 3 {
			finalTriades = append(finalTriades, tempTriades[x])
		}
	}

	return finalTriades
}

func buildTriadeSymbols(exchangeInfo models.ExchangeInfo, triade models.Triade) []models.Symbol {
	symbols := []models.Symbol{}
	for c, symbol := range triade.Assets {
		a := symbol
		b := triade.Assets[(c+1)%3]
		order1 := a + b
		order2 := b + a
		for _, pair := range exchangeInfo.Symbols {
			if pair.Status != "TRADING" {
				continue
			}
			if pair.Symbol == order1 {
				symbols = append(symbols, pair)
			} else if pair.Symbol == order2 {
				symbols = append(symbols, pair)
			}
		}
	}
	return symbols
}

func BuildStreamURL(pairSymbols []string) string {
	url := "wss://stream.binance.com:9443/ws/"
	for i, symbol := range pairSymbols {
		url += strings.ToLower(symbol) + "@trade"
		if i < len(pairSymbols)-1 {
			url += "/"
		}
	}
	return url
}
