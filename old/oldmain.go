package old

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Trade struct {
	EventType string `json:"e"`
	EventTime int    `json:"E"`
	Symbol    string `json:"s"`
	TradeID   int    `json:"t"`
	Price     string `json:"p"`
	Quantity  string `json:"q"`
	TradeTime int    `json:"T"`
	Maker     bool   `json:"m"`
	BestMatch bool   `json:"M"`
}

var (
	priceMap = struct {
		sync.RWMutex
		m map[string]string
	}{m: make(map[string]string)}
)

var USDTBalance float64 = 100
var orderRequestCounter int = 0

func main() {
	fmt.Println("Starting Arbitrage bot V0!")

	// Paires de tokens à écouter
	pairs := []string{"btcusdt", "ethbtc", "ethusdt"}

	var wg sync.WaitGroup

	for _, pair := range pairs {
		wg.Add(1)
		go func(pair string) {
			defer wg.Done()
			listenForTrades(pair)
		}(pair)
	}

	wg.Wait()
}

func listenForTrades(pair string) {
	// create a new websocket connection
	url := "wss://stream.binance.com:9443/ws/" + pair + "@trade"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal(err, url)
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}

		var trade Trade
		err = json.Unmarshal(message, &trade)
		if err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		priceMap.Lock()
		priceMap.m[trade.Symbol] = trade.Price
		priceMap.Unlock()

		// Check for arbitrage opportunities after updating the price map
		checkArbitrage()
	}
}

func checkArbitrage() {
	priceMap.RLock()
	defer priceMap.RUnlock()

	btcPriceInUsdt := priceMap.m["BTCUSDT"]

	btcPriceInUsdtFloat, err := strconv.ParseFloat(btcPriceInUsdt, 32)
	if err != nil {

		return
	}

	ethPriceInBtc := priceMap.m["ETHBTC"]

	ethPriceInBtcFloat, err := strconv.ParseFloat(ethPriceInBtc, 32)
	if err != nil {
		return
	}

	ethPriceInUsdt := priceMap.m["ETHUSDT"]

	ethPriceInUsdtFloat, err := strconv.ParseFloat(ethPriceInUsdt, 32)
	if err != nil {
		return
	}

	usdtPriceInEthFloat := 1 / ethPriceInUsdtFloat
	// convert eth price in usdt to usdt price in eth

	log.Println("BTC price: ", btcPriceInUsdtFloat, " usdt")
	log.Println("ETH price: ", ethPriceInBtcFloat, " btc")
	log.Println("USDT price: ", usdtPriceInEthFloat, " eth")

	fmt.Println("Simulating arbitrage with " + fmt.Sprint(USDTBalance) + " USDT")
	virtualBalance := USDTBalance
	previousUSDTBalance := virtualBalance

	btcAmount := virtualBalance / btcPriceInUsdtFloat
	virtualBalance = 0
	fmt.Println("Buying " + fmt.Sprint(btcAmount) + " btc with usdt")

	ethAmount := btcAmount / ethPriceInBtcFloat
	btcAmount = 0
	fmt.Println("Buying " + fmt.Sprint(ethAmount) + " eth with btc")

	virtualBalance = ethAmount / usdtPriceInEthFloat
	ethAmount = 0
	fmt.Println("Buying" + fmt.Sprint(virtualBalance) + " usdt with eth")
	// profit percentage
	profit := (virtualBalance - previousUSDTBalance) / previousUSDTBalance * 100
	fmt.Println("Profit percentage: ", profit)
	fmt.Println("Order request counter: ", orderRequestCounter)
	if profit > 0 {
		orderRequestCounter++
		if orderRequestCounter == 100 {
			fmt.Println("Max order request reached, not taking the opportunity and waiting for 10 seconds to reset the counter.")
			// goroutine that prints a dot every second
			go func() {
				for {
					fmt.Print(".")
					time.Sleep(1 * time.Second)
				}
			}()
			time.Sleep(10 * time.Second)
			orderRequestCounter = 0
			return
		}

		USDTBalance = virtualBalance
		fmt.Println("Arbitrage opportunity found, new balance: ", USDTBalance)
	}
	fmt.Println("-------------------------------------------------")
	fmt.Println("actual balance: ", USDTBalance)
}
