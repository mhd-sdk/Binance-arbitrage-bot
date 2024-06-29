package instance

import (
	"DeltA/pkg/models"
	"fmt"
	"math/big"

	binance_connector "github.com/binance/binance-connector-go"
)

func FindPrice(from string, to string, symbolPrices map[string]*big.Float) *big.Float {
	_, ok := symbolPrices[to+from]
	if ok {
		return symbolPrices[to+from]
	} else {
		// return 1 / symbolPrices[from+to]
		// divide by 1 to get the inverse price
		return new(big.Float).Quo(big.NewFloat(1), symbolPrices[from+to])
	}
}

type OrderPairs []OrderPair

type OrderPair struct {
	From   string
	To     string
	Price  *big.Float
	Symbol *binance_connector.SymbolInfo
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

		i.Mutex.Lock()
		price := FindPrice(from, to, i.SymbolPrices)
		i.Mutex.Unlock()
		symbol := models.FindSymbol(from, to, triade.Symbols)

		orderedPairs1 = append(orderedPairs1, OrderPair{
			From:   from,
			To:     to,
			Price:  price,
			Symbol: symbol,
		})
	}

	// orderedSymbols1, _ := triade.OrderSymbols(i.StartingStableAsset)

	for _, pair := range orderedPairs1 {
		calculatedBalance[pair.To] = new(big.Float).Quo(calculatedBalance[pair.From], pair.Price)
		calculatedBalance[pair.From] = big.NewFloat(0)
	}

	diffBetweenBalances := new(big.Float).Sub(calculatedBalance[i.StartingStableAsset], initialBalance[i.StartingStableAsset])

	gainPercentage := new(big.Float).Quo(new(big.Float).Mul(diffBetweenBalances, big.NewFloat(100)), initialBalance[i.StartingStableAsset])
	fmt.Println(gainPercentage)
	if gainPercentage.Cmp(big.NewFloat(1)) == 1 {
		logs := "Triangular arbitrage opportunity found : " + order1[0] + "/" + order1[1] + "/" + order1[2] + "Gain: " + gainPercentage.Text('f', -1) + "%"
		i.Logger.Slog.Info(logs)

		filePath := "logs/" + triade.Assets[0] + "-" + triade.Assets[1] + "-" + triade.Assets[2] + ".txt"
		i.Logger.FileLog(filePath, logs)
		i.ComputeTriangularOrders(orderedPairs1)
	} else if gainPercentage.Cmp(big.NewFloat(-1)) == -1 {
		// reverse the order

		var initialBalance = Balance{
			i.StartingStableAsset: big.NewFloat(100),
		}
		var calculatedBalance = Balance{
			i.StartingStableAsset: big.NewFloat(100),
		}

		// find the two order starting from the first asset
		_, order2 := triade.OrderAssets(i.StartingStableAsset)

		orderedPairs2 := OrderPairs{}

		for c, asset := range order2 {
			from := asset
			to := order2[(c+1)%3]

			i.Mutex.Lock()
			price := FindPrice(from, to, i.SymbolPrices)
			i.Mutex.Unlock()
			symbol := models.FindSymbol(from, to, triade.Symbols)

			orderedPairs2 = append(orderedPairs2, OrderPair{
				From:   from,
				To:     to,
				Price:  price,
				Symbol: symbol,
			})
		}

		// orderedSymbols1, _ := triade.OrderSymbols(i.StartingStableAsset)

		for _, pair := range orderedPairs2 {
			calculatedBalance[pair.To] = new(big.Float).Quo(calculatedBalance[pair.From], pair.Price)
			calculatedBalance[pair.From] = big.NewFloat(0)
		}

		diffBetweenBalances := new(big.Float).Sub(calculatedBalance[i.StartingStableAsset], initialBalance[i.StartingStableAsset])

		gainPercentage := new(big.Float).Quo(new(big.Float).Mul(diffBetweenBalances, big.NewFloat(100)), initialBalance[i.StartingStableAsset])
		if gainPercentage.Cmp(big.NewFloat(1)) == 1 {
			logs := "Triangular arbitrage opportunity found : " + order2[0] + "/" + order2[1] + "/" + order2[2] + "Gain: " + gainPercentage.Text('f', -1) + "%"
			i.Logger.Slog.Info(logs)

			filePath := "logs/" + triade.Assets[0] + "-" + triade.Assets[1] + "-" + triade.Assets[2] + ".txt"
			i.Logger.FileLog(filePath, logs)
			i.ComputeTriangularOrders(orderedPairs2)
		}
	}
}
