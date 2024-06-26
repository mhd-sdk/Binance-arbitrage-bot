package helpers

import (
	"DeltA/pkg/binance"
	"DeltA/pkg/models"
	"slices"
)

func ExInfoToAvailableAssetsEx(exchangeInfo binance.ExchangeInfo) []string {
	availableAssets := []string{}
	for _, pair := range exchangeInfo {
		if !slices.Contains(availableAssets, pair.FromAsset) {
			availableAssets = append(availableAssets, pair.FromAsset)
		}
		if !slices.Contains(availableAssets, pair.ToAsset) {
			availableAssets = append(availableAssets, pair.ToAsset)
		}
	}
	return availableAssets
}

func checkIsExchangeAvailable(assets [3]string, exchangeInfo binance.ExchangeInfo) bool {
	for _, pair := range exchangeInfo {
		if (pair.FromAsset == assets[0] && pair.ToAsset == assets[1]) || (pair.FromAsset == assets[1] && pair.ToAsset == assets[0]) {
			return true
		}
		if (pair.FromAsset == assets[1] && pair.ToAsset == assets[2]) || (pair.FromAsset == assets[2] && pair.ToAsset == assets[1]) {
			return true
		}
		if (pair.FromAsset == assets[0] && pair.ToAsset == assets[2]) || (pair.FromAsset == assets[2] && pair.ToAsset == assets[0]) {
			return true
		}
	}
	return false
}

func BuildTriades(exchangeInfo binance.ExchangeInfo) []models.Triade {
	availableAssetsExchange := ExInfoToAvailableAssetsEx(exchangeInfo)
	triades := []models.Triade{}
	for i := 0; i < len(availableAssetsExchange); i++ {
		for j := i + 1; j < len(availableAssetsExchange); j++ {
			for k := j + 1; k < len(availableAssetsExchange); k++ {
				if availableAssetsExchange[i] != "USDT" && availableAssetsExchange[j] != "USDT" && availableAssetsExchange[k] != "USDT" {
					continue
				}
				assets := [3]string{availableAssetsExchange[i], availableAssetsExchange[j], availableAssetsExchange[k]}

				pairs := []models.Pair{}
				for _, pair := range exchangeInfo {
					if (pair.FromAsset == assets[0] && pair.ToAsset == assets[1]) && pair.FromIsBase {
						pairs = append(pairs, pair)
					}
					if (pair.FromAsset == assets[1] && pair.ToAsset == assets[2]) && pair.FromIsBase {
						pairs = append(pairs, pair)
					}
					if (pair.FromAsset == assets[0] && pair.ToAsset == assets[2]) && pair.FromIsBase {
						pairs = append(pairs, pair)
					}
				}
				if len(pairs) != 3 {
					continue
				}
				triades = append(triades, models.Triade{
					Assets: [3]string{availableAssetsExchange[i], availableAssetsExchange[j], availableAssetsExchange[k]},
					Pairs:  pairs,
				})
			}
		}
	}
	return triades
}

func FormatTriades(triades []models.Triade) string {
	formatted := ""
	for _, triade := range triades {
		formatted += triade.Assets[0] + " -> " + triade.Assets[1] + " -> " + triade.Assets[2] + "\n"
	}
	return formatted
}
