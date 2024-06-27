package instance

import (
	"DeltA/pkg/models"
	"fmt"
)

func FindPrice(from string, to string, symbolPrices map[string]float64) float64 {
	_, ok := symbolPrices[to+from]
	if ok {
		return symbolPrices[to+from]
	} else {
		return 1 / symbolPrices[from+to]
	}
}

type OrderedPair []Pair

type Pair struct {
	From                   string
	To                     string
	StringifiedTransaction string
	Price                  float64
}

type Balance map[string]float64

var balance = Balance{
	"USDT": 100.0,
}

func (i *Instance) ScanOpportunities(triade models.Triade) (float64, string) {

	var virtualBalance = make(Balance)
	for key, value := range balance {
		virtualBalance[key] = value
	}

	// find the two order starting from the first asset
	order1, _ := triade.OrderAssets(i.StartingStableAsset)

	orderedPairs := OrderedPair{}

	for c, asset := range order1 {
		from := asset
		to := order1[(c+1)%3]
		orderedPairs = append(orderedPairs, Pair{
			StringifiedTransaction: from + "->" + to,
			From:                   from,
			To:                     to,
			Price:                  FindPrice(from, to, i.SymbolPrices),
		})
	}

	stringifiedTransactions := ""
	for _, pair := range orderedPairs {
		stringifiedTransactions += pair.StringifiedTransaction + " at " + fmt.Sprint(pair.Price)
	}

	var lastBalance = make(Balance)
	for key, value := range virtualBalance {
		lastBalance[key] = value
	}

	for _, pair := range orderedPairs {
		virtualBalance[pair.To] = virtualBalance[pair.From] / pair.Price
		virtualBalance[pair.From] = 0
	}

	lastBalance[i.StartingStableAsset] = balance[i.StartingStableAsset]

	// gain := strconv.FormatFloat(virtualBalance[i.StartingStableAsset]-lastBalance[i.StartingStableAsset], 'g', -1, 64)

	// gainMessage := "Gained: " + "("
	// if virtualBalance[i.StartingStableAsset] > lastBalance[i.StartingStableAsset] {
	// 	gainMessage += "+"
	// }
	// gainMessage += gain + ")"

	// i.Logger.Info(gainMessage)
	gainPercentage := (virtualBalance[i.StartingStableAsset] - lastBalance[i.StartingStableAsset]) / lastBalance[i.StartingStableAsset]
	if gainPercentage > 0.1 || gainPercentage < -0.1 {
		i.Logger.Info("Gain % : " + fmt.Sprint(gainPercentage) + "%")
		fmt.Println("-------------------")
	}

	// i.Logger.Info("New balance: " + fmt.Sprint(virtualBalance) + "  (+" + gainString + orderedPairs[0].From + ")")

	// gainPercentage := virtualBalance[triade.Assets[0]] / lastBalance[triade.Assets[0]] * 100
	// fmt.Println("Gain % :", gainPercentage, "%")
	return 0, ""

	// old manual way

	// btcPriceInUsdt := priceMap.m["BTCUSDT"]

	// btcPriceInUsdtFloat, err := strconv.ParseFloat(btcPriceInUsdt, 32)
	// if err != nil {

	// 	return
	// }

	// ethPriceInBtc := priceMap.m["ETHBTC"]

	// ethPriceInBtcFloat, err := strconv.ParseFloat(ethPriceInBtc, 32)
	// if err != nil {
	// 	return
	// }

	// ethPriceInUsdt := priceMap.m["ETHUSDT"]

	// ethPriceInUsdtFloat, err := strconv.ParseFloat(ethPriceInUsdt, 32)
	// if err != nil {
	// 	return
	// }

	// usdtPriceInEthFloat := 1 / ethPriceInUsdtFloat
	// // convert eth price in usdt to usdt price in eth

	// // clear the console
	// fmt.Print("\033[H\033[2J")

	// log.Println("BTC price: ", btcPriceInUsdtFloat, " usdt")
	// log.Println("ETH price: ", ethPriceInBtcFloat, " btc")
	// log.Println("USDT price: ", usdtPriceInEthFloat, " eth")

	// fmt.Println("Simulating arbitrage with " + fmt.Sprint(USDTBalance) + " USDT")
	// virtualBalance := USDTBalance
	// previousUSDTBalance := virtualBalance

	// btcAmount := virtualBalance / btcPriceInUsdtFloat
	// virtualBalance = 0
	// fmt.Println("Buying " + fmt.Sprint(btcAmount) + " btc with usdt")

	// ethAmount := btcAmount / ethPriceInBtcFloat
	// btcAmount = 0
	// fmt.Println("Buying " + fmt.Sprint(ethAmount) + " eth with btc")

	// virtualBalance = ethAmount / usdtPriceInEthFloat
	// ethAmount = 0
	// fmt.Println("Buying" + fmt.Sprint(virtualBalance) + " usdt with eth")
	// // profit percentage
	// profit := (virtualBalance - previousUSDTBalance) / previousUSDTBalance * 100
	// fmt.Println("Profit percentage: ", profit)
	// fmt.Println("Order request counter: ", orderRequestCounter)
	// if profit > 0 {
	// 	orderRequestCounter++
	// 	if orderRequestCounter == 100 {
	// 		fmt.Println("Max order request reached, not taking the opportunity and waiting for 10 seconds to reset the counter.")
	// 		// goroutine that prints a dot every second
	// 		go func() {
	// 			for {
	// 				fmt.Print(".")
	// 				time.Sleep(1 * time.Second)
	// 			}
	// 		}()
	// 		time.Sleep(10 * time.Second)
	// 		orderRequestCounter = 0
	// 		return
	// 	}

	// 	USDTBalance = virtualBalance
	// 	fmt.Println("Arbitrage opportunity found, new balance: ", USDTBalance)
	// }
	// fmt.Println("-------------------------------------------------")
	// fmt.Println("actual balance: ", USDTBalance)
}
