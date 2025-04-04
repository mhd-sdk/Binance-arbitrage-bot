package helpers

import (
	"DeltA/pkg/models"
	"math"
	"slices"
	"strings"

	binance_connector "github.com/binance/binance-connector-go"
)

func ListAllAssets(exchangeInfo *binance_connector.ExchangeInfoResponse) []string {
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

func BuildTriades(exchangeInfo *binance_connector.ExchangeInfoResponse, startingAsset string) []models.Triade {
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

func buildTriadeSymbols(exchangeInfo *binance_connector.ExchangeInfoResponse, triade models.Triade) []*binance_connector.SymbolInfo {
	symbols := []*binance_connector.SymbolInfo{}
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

func FloorToPrecision(num float64, precision int) float64 {
	factor := math.Pow(10, float64(precision))
	return math.Floor(num*factor) / factor
}

func CheckAssetIsBase(asset string, symbol *binance_connector.SymbolInfo) bool {
	return asset == symbol.BaseAsset
}

func GetSymbolInfo(symbol string, exchangeInfo *binance_connector.ExchangeInfoResponse) *binance_connector.SymbolInfo {
	for _, sym := range exchangeInfo.Symbols {
		if sym.Symbol == symbol {
			return sym
		}
	}
	return nil
}
