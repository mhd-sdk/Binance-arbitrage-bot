package instance

import (
	"DeltA/pkg/models"
	"fmt"
	"math/big"
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

type OrderedPair []Pair

type Pair struct {
	From  string
	To    string
	Price *big.Float
}

type Balance map[string]*big.Float

func (i Instance) ScanOpportunities(triade models.Triade) {
	var initialBalance = Balance{
		"USDT": big.NewFloat(100),
	}
	var calculatedBalance = Balance{
		"USDT": big.NewFloat(100),
	}

	// find the two order starting from the first asset
	order1, _ := triade.OrderAssets(i.StartingStableAsset)

	orderedPairs := OrderedPair{}

	for c, asset := range order1 {
		from := asset
		to := order1[(c+1)%3]

		i.Mutex.Lock()
		price := FindPrice(from, to, i.SymbolPrices)
		i.Mutex.Unlock()
		orderedPairs = append(orderedPairs, Pair{
			From:  from,
			To:    to,
			Price: price,
		})
	}

	fmt.Println("-----------------")

	for _, pair := range orderedPairs {
		calculatedBalance[pair.To] = new(big.Float).Quo(calculatedBalance[pair.From], pair.Price)
		calculatedBalance[pair.From] = big.NewFloat(0)
	}

	diffBetweenBalances := new(big.Float).Sub(calculatedBalance[i.StartingStableAsset], initialBalance[i.StartingStableAsset])

	i.Logger.Info("Initial balance : " + fmt.Sprint(initialBalance[i.StartingStableAsset]) + i.StartingStableAsset)
	i.Logger.Info("New balance : " + fmt.Sprint(calculatedBalance[i.StartingStableAsset]) + i.StartingStableAsset)
	i.Logger.Info("Difference : " + diffBetweenBalances.Text('f', -1) + i.StartingStableAsset)

	gainPercentage := new(big.Float).Quo(new(big.Float).Mul(diffBetweenBalances, big.NewFloat(100)), initialBalance[i.StartingStableAsset])

	i.Logger.Info("Gain: : " + diffBetweenBalances.Text('f', -1) + " " + i.StartingStableAsset + " (" + gainPercentage.Text('f', -1) + "%)")
	// printf l'expression

	// if gainPercentage.Cmp(big.NewFloat(0.3)) == 1 {
	// 	i.Logger.Info("Arbitrage opportunity found: ")
	// 	i.Logger.Info("Gain % : " + fmt.Sprint(gainPercentage) + "%")
	// 	i.Logger.Info("New balance: " + fmt.Sprint(calculatedBalance[i.StartingStableAsset]) + "  (+" + orderedPairs[0].From + ")")
	// 	fmt.Println("---------A----------")
	// }
	// else if gainPercentage.Cmp(big.NewFloat(-0.3)) == -1 {

	// 	var virtualBalance = Balance{
	// 		"USDT": big.NewFloat(100),
	// 	}

	// 	// find the two order starting from the first asset
	// 	_, order2 := triade.OrderAssets(i.StartingStableAsset)

	// 	orderedPairs := OrderedPair{}

	// 	for c, asset := range order2 {
	// 		from := asset
	// 		to := order2[(c+1)%3]
	// 		orderedPairs = append(orderedPairs, Pair{
	// 			From:  from,
	// 			To:    to,
	// 			Price: FindPrice(from, to, i.SymbolPrices),
	// 		})
	// 	}

	// 	for _, pair := range orderedPairs {
	// 		virtualBalance[pair.To] = virtualBalance[pair.From] / pair.Price
	// 		virtualBalance[pair.From] = 0
	// 	}

	// 	gainPercentage := (virtualBalance[i.StartingStableAsset] - balance[i.StartingStableAsset]) / balance[i.StartingStableAsset]
	// 	i.Logger.Info("Gain % : " + fmt.Sprint(gainPercentage) + "%")
	// 	i.Logger.Info("StartingBalance: " + fmt.Sprint(balance[i.StartingStableAsset]))
	// 	i.Logger.Info("New balance: " + fmt.Sprint(virtualBalance[i.StartingStableAsset]))
	// 	fmt.Println("---------V----------")
	// }

}
