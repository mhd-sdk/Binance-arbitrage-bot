package instance

import (
	"DeltA/pkg/binance"
	"DeltA/pkg/helpers"
	"DeltA/pkg/logging"
	"DeltA/pkg/models"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Instance struct {
	Mutex               *sync.Mutex
	ExchangeInfo        models.ExchangeInfo
	SymbolPrices        map[string]*big.Float
	StartingStableAsset string
	Triades             []models.Triade
	Logger              *logging.Logger
}

func NewInstance() *Instance {
	return &Instance{}
}

func (i *Instance) Init() (err error) {
	i.Mutex = new(sync.Mutex)
	logger, err := logging.NewLogger()
	i.Logger = logger
	i.Logger.Slog.Info("Starting DeltΔ  BOT")
	if err != nil {
		logger.Slog.Error(err.Error())
		return err
	}

	i.StartingStableAsset = "USDC"
	i.Logger.Slog.Info("Loading exchange info from Binance...")
	i.ExchangeInfo, err = binance.GetExchangeInfo()
	if err != nil {
		return err
	}

	i.SymbolPrices, err = binance.GetSymbolPrices()
	if err != nil {
		return err
	}

	i.Triades = helpers.BuildTriades(i.ExchangeInfo, i.StartingStableAsset)

	i.Logger.Slog.Info(fmt.Sprintf("Generated %d triades", len(i.Triades)))
	return nil
}

func (i *Instance) Watch() {
	var wg sync.WaitGroup
	wg.Add(len(i.Triades))
	// wg.Add(200)

	for _, triade := range i.Triades {
		go func(triade models.Triade) {
			defer wg.Done()
			pairSymbols := make([]string, len(triade.Symbols))
			for idx, symbol := range triade.Symbols {
				pairSymbols[idx] = symbol.Symbol
			}

			url := helpers.BuildStreamURL(pairSymbols)
			i.Logger.Slog.Info(fmt.Sprintf("Connecting to %s", url))

			for {

				conn, _, err := websocket.DefaultDialer.Dial(url, nil)
				if err != nil {
					i.Logger.Slog.Error(fmt.Sprintf("Error connecting to %s: %s", url, err))
					time.Sleep(5 * time.Second) // Retry after 5 seconds
					continue
				}

				defer conn.Close()
				i.Logger.Slog.Info(fmt.Sprintf("Watching for arbitrage opportunities... (%s/%s/%s)", triade.Assets[0], triade.Assets[1], triade.Assets[2]))

				for {
					_, message, err := conn.ReadMessage()
					if err != nil {
						i.Logger.Slog.Error("Read error:", err)
						break
					}

					var trade models.Trade
					err = json.Unmarshal(message, &trade)
					if err != nil {
						i.Logger.Slog.Error("Unmarshal error:", err)
						continue
					}

					// price, err := strconv.ParseFloat(trade.Price, 64)
					// parse big float
					price, _, err := big.ParseFloat(trade.Price, 10, 64, big.ToNearestEven)
					if err != nil {
						i.Logger.Slog.Error(err.Error())
						continue
					}
					i.Mutex.Lock()
					i.SymbolPrices[trade.Symbol] = price
					i.Mutex.Unlock()
					// if i.SymbolPrices[triade.Assets[0]] == nil || i.SymbolPrices[triade.Assets[1]] == nil || i.SymbolPrices[triade.Assets[2]] == nil {
					// 	// print missing prices
					// 	i.Logger.Slog.Info(fmt.Sprintf("Missing prices for %s/%s/%s", triade.Assets[0], triade.Assets[1], triade.Assets[2]))
					// 	continue
					// }

					i.ScanOpportunities(triade)
				}
			}
		}(triade)
	}
	wg.Wait()
}
