package instance

import (
	"DeltA/pkg/helpers"
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
		i.StartingStableAsset: big.NewFloat(100),
	}
	var calculatedBalance = Balance{
		i.StartingStableAsset: big.NewFloat(100),
	}

	// find the two order starting from the first asset
	order1, _ := triade.OrderAssets(i.StartingStableAsset)

	orderedPairs1 := OrderedPair{}

	for c, asset := range order1 {
		from := asset
		to := order1[(c+1)%3]

		i.Mutex.Lock()
		price := FindPrice(from, to, i.SymbolPrices)
		i.Mutex.Unlock()

		orderedPairs1 = append(orderedPairs1, Pair{
			From:  from,
			To:    to,
			Price: price,
		})
	}

	orderedSymbols1, _ := triade.OrderSymbols(i.StartingStableAsset)

	for _, pair := range orderedPairs1 {
		calculatedBalance[pair.To] = new(big.Float).Quo(calculatedBalance[pair.From], pair.Price)
		calculatedBalance[pair.From] = big.NewFloat(0)
	}

	diffBetweenBalances := new(big.Float).Sub(calculatedBalance[i.StartingStableAsset], initialBalance[i.StartingStableAsset])

	gainPercentage := new(big.Float).Quo(new(big.Float).Mul(diffBetweenBalances, big.NewFloat(100)), initialBalance[i.StartingStableAsset])

	// printf l'expression
	if gainPercentage.Cmp(big.NewFloat(1.5)) == 1 {
		log1 := "Opportunity found for triade: " + order1[0] + "/" + order1[1] + "/" + order1[2]
		log0 := helpers.BuildStreamURL([]string{order1[0], order1[1], order1[2]})
		log2 := "Initial balance : " + fmt.Sprint(initialBalance[i.StartingStableAsset]) + i.StartingStableAsset
		log3 := "New balance : " + fmt.Sprint(calculatedBalance[i.StartingStableAsset]) + i.StartingStableAsset
		log4 := "Gain: : " + diffBetweenBalances.Text('f', -1) + " " + i.StartingStableAsset + " (" + gainPercentage.Text('f', -1) + "%)"
		log5 := orderedSymbols1[0].Symbol + "\n" + orderedSymbols1[1].Symbol + "\n" + orderedSymbols1[2].Symbol
		logs := log1 + "\n" + log0 + "\n" + log2 + "\n" + log3 + "\n" + log4 + "\n" + fmt.Sprint(log5)
		i.Logger.Slog.Info(logs)

		filePath := "logs/" + triade.Assets[0] + "-" + triade.Assets[1] + "-" + triade.Assets[2] + ".txt"
		i.Logger.FileLog(filePath, logs)
	}
	// else {

	// 	var initialBalance = Balance{
	// 		i.StartingStableAsset: big.NewFloat(100),
	// 	}
	// 	var calculatedBalance = Balance{
	// 		i.StartingStableAsset: big.NewFloat(100),
	// 	}

	// 	orderedPairs2 := OrderedPair{}

	// 	for c, asset := range order2 {
	// 		from := asset
	// 		to := order2[(c+1)%3]

	// 		i.Mutex.Lock()
	// 		price := FindPrice(from, to, i.SymbolPrices)
	// 		i.Mutex.Unlock()

	// 		orderedPairs2 = append(orderedPairs2, Pair{
	// 			From:  from,
	// 			To:    to,
	// 			Price: price,
	// 		})
	// 	}

	// 	for _, pair := range orderedPairs2 {
	// 		calculatedBalance[pair.To] = new(big.Float).Quo(calculatedBalance[pair.From], pair.Price)
	// 		calculatedBalance[pair.From] = big.NewFloat(0)
	// 	}

	// 	diffBetweenBalances := new(big.Float).Sub(calculatedBalance[i.StartingStableAsset], initialBalance[i.StartingStableAsset])

	// 	gainPercentage := new(big.Float).Quo(new(big.Float).Mul(diffBetweenBalances, big.NewFloat(100)), initialBalance[i.StartingStableAsset])

	// 	if gainPercentage.Cmp(big.NewFloat(2)) == 1 {
	// 		log1 := "Opportunity found for triade: " + order2[0] + "/" + order2[1] + "/" + order2[2]
	// 		log2 := "Initial balance : " + fmt.Sprint(initialBalance[i.StartingStableAsset]) + i.StartingStableAsset
	// 		log3 := "New balance : " + fmt.Sprint(calculatedBalance[i.StartingStableAsset]) + i.StartingStableAsset
	// 		log4 := "Gain: : " + diffBetweenBalances.Text('f', -1) + " " + i.StartingStableAsset + " (" + gainPercentage.Text('f', -1) + "%)"
	// 		log5 := orderedSymbols2
	// 		logs := log1 + "\n" + log2 + "\n" + log3 + "\n" + log4 + "\n" + fmt.Sprint(log5)
	// 		i.Logger.Slog.Info(logs)
	// 	}
	// }

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
