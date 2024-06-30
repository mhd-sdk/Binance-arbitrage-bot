package instance

import (
	"DeltA/pkg/helpers"
	"DeltA/pkg/models"
	"math/big"

	binance_connector "github.com/binance/binance-connector-go"
)

type OrderPairs []OrderPair

type OrderPair struct {
	From        string
	To          string
	SymbolPrice *big.Float
	Symbol      *binance_connector.SymbolInfo
}

type Balance map[string]*big.Float

func (i Instance) ScanOpportunities(triade models.Triade) {
	var initialBalance = Balance{
		i.StartingStableAsset: big.NewFloat(100),
	}
	var calculatedBalance = Balance{
		i.StartingStableAsset: big.NewFloat(100),
	}

	// find the two order starting from the first asset
	order1, _ := triade.OrderAssets(i.StartingStableAsset)

	orderedPairs1 := OrderPairs{}

	for c, asset := range order1 {
		from := asset
		to := order1[(c+1)%3]

		symbol := models.FindSymbol(from, to, triade.Symbols)

		i.Mutex.Lock()
		price := i.SymbolPrices[symbol.Symbol]
		i.Mutex.Unlock()

		orderedPairs1 = append(orderedPairs1, OrderPair{
			From:        from,
			To:          to,
			SymbolPrice: price,
			Symbol:      symbol,
		})
	}

	for _, pair := range orderedPairs1 {
		if helpers.CheckAssetIsBase(pair.To, pair.Symbol) {
			calculatedBalance[pair.To] = new(big.Float).Quo(calculatedBalance[pair.From], pair.SymbolPrice)
			calculatedBalance[pair.From] = big.NewFloat(0)
		} else {
			calculatedBalance[pair.To] = new(big.Float).Mul(calculatedBalance[pair.From], pair.SymbolPrice)
			calculatedBalance[pair.From] = big.NewFloat(0)
		}
	}

	gainPercentage := new(big.Float).Sub(calculatedBalance[i.StartingStableAsset], initialBalance[i.StartingStableAsset])

	if gainPercentage.Cmp(big.NewFloat(2)) == 1 && gainPercentage.Cmp(big.NewFloat(20)) == -1 {
		// logs := "Triangular arbitrage opportunity found : " + order1[0] + "/" + order1[1] + "/" + order1[2] + " Gain: " + gainPercentage.Text('f', -1) + "%"
		// i.Logger.Slog.Info(logs)

		// logs = "Binance URLS: " + binanceurl1 + " " + binanceurl2 + " " + binanceurl3
		// i.Logger.Slog.Info(logs)

		// filePath := "logs/" + triade.Assets[0] + "-" + triade.Assets[1] + "-" + triade.Assets[2] + ".txt"
		// i.Logger.FileLog(filePath, logs)
		i.ComputeTriangularOrders(orderedPairs1)
	}
	// else if gainPercentage.Cmp(big.NewFloat(-0)) == -1 {
	// 	// reverse the order

	// 	var initialBalance = Balance{
	// 		i.StartingStableAsset: big.NewFloat(100),
	// 	}
	// 	var calculatedBalance = Balance{
	// 		i.StartingStableAsset: big.NewFloat(100),
	// 	}

	// 	// find the two order starting from the first asset
	// 	_, order2 := triade.OrderAssets(i.StartingStableAsset)

	// 	orderedPairs2 := OrderPairs{}

	// 	for c, asset := range order2 {
	// 		from := asset
	// 		to := order2[(c+1)%3]

	// 		i.Mutex.Lock()
	// 		basePrice, price := FindPrice(from, to, i.SymbolPrices)
	// 		i.Mutex.Unlock()
	// 		symbol := models.FindSymbol(from, to, triade.Symbols)

	// 		orderedPairs2 = append(orderedPairs2, OrderPair{
	// 			From:        from,
	// 			To:          to,
	// 			Price:       price,
	// 			SymbolPrice: basePrice,
	// 			Symbol:      symbol,
	// 		})
	// 	}

	// 	for _, pair := range orderedPairs2 {
	// 		calculatedBalance[pair.To] = new(big.Float).Mul(new(big.Float).Quo(calculatedBalance[pair.From], pair.Price), big.NewFloat(0.9))
	// 		calculatedBalance[pair.From] = big.NewFloat(0)
	// 	}

	// 	diffBetweenBalances := new(big.Float).Sub(calculatedBalance[i.StartingStableAsset], initialBalance[i.StartingStableAsset])

	// 	gainPercentage := new(big.Float).Quo(new(big.Float).Mul(diffBetweenBalances, big.NewFloat(100)), initialBalance[i.StartingStableAsset])
	// 	i.Logger.Slog.Info("Gain percentage: %s for %s", gainPercentage.Text('f', -1), order1)
	// 	if gainPercentage.Cmp(big.NewFloat(3)) == 1 && gainPercentage.Cmp(big.NewFloat(20)) == -1 {
	// 		logs := "Triangular arbitrage opportunity found : " + order2[0] + "/" + order2[1] + "/" + order2[2] + " Gain: " + gainPercentage.Text('f', -1) + "%"
	// 		i.Logger.Slog.Info(logs)

	// 		// binanceurl1 := "https://www.binance.com/en/trade/" + orderedPairs1[0].Symbol.Symbol
	// 		// binanceurl2 := "https://www.binance.com/en/trade/" + orderedPairs1[1].Symbol.Symbol
	// 		// binanceurl3 := "https://www.binance.com/en/trade/" + orderedPairs1[2].Symbol.Symbol
	// 		// logs = "Binance URLS: " + binanceurl1 + " " + binanceurl2 + " " + binanceurl3
	// 		// i.Logger.Slog.Info(logs)

	// 		// filePath := "logs/" + triade.Assets[0] + "-" + triade.Assets[1] + "-" + triade.Assets[2] + ".txt"
	// 		// i.Logger.FileLog(filePath, logs)
	// 		// i.ComputeTriangularOrders(orderedPairs2)
	// 	}
	// }
}
