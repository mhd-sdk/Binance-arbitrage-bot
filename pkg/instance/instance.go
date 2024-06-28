package instance

import (
	"DeltA/pkg/binance"
	"DeltA/pkg/helpers"
	"DeltA/pkg/models"
	"encoding/json"
	"fmt"
	"log/slog"
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
	Logger              *slog.Logger
}

func NewInstance(logger *slog.Logger) *Instance {
	mu := new(sync.Mutex)
	return &Instance{Logger: logger, Mutex: mu}
}

func (i *Instance) Init() (err error) {
	i.StartingStableAsset = "USDT"
	i.Logger.Info("Loading exchange info from Binance...")
	i.ExchangeInfo, err = binance.GetExchangeInfo()
	if err != nil {
		return err
	}

	i.SymbolPrices, err = binance.GetSymbolPrices()
	if err != nil {
		return err
	}

	i.Triades = helpers.BuildTriades(i.ExchangeInfo, i.StartingStableAsset)

	i.Logger.Info(fmt.Sprintf("Generated %d triades", len(i.Triades)))
	return nil
}

func (i *Instance) Watch() {
	var wg sync.WaitGroup
	// wg.Add(len(i.Triades))
	wg.Add(200)

	for _, triade := range i.Triades[:200] {
		go func(triade models.Triade) {
			defer wg.Done()
			pairSymbols := make([]string, len(triade.Symbols))
			for idx, symbol := range triade.Symbols {
				pairSymbols[idx] = symbol.Symbol
			}

			url := helpers.BuildStreamURL(pairSymbols)
			i.Logger.Info(fmt.Sprintf("Connecting to %s", url))

			for {

				conn, _, err := websocket.DefaultDialer.Dial(url, nil)
				if err != nil {
					i.Logger.Error(fmt.Sprintf("Error connecting to %s: %s", url, err))
					time.Sleep(5 * time.Second) // Retry after 5 seconds
					continue
				}

				defer conn.Close()
				i.Logger.Info(fmt.Sprintf("Watching for arbitrage opportunities... (%s/%s/%s)", triade.Assets[0], triade.Assets[1], triade.Assets[2]))

				for {
					_, message, err := conn.ReadMessage()
					if err != nil {
						i.Logger.Error("Read error:", err)
						break
					}

					var trade models.Trade
					err = json.Unmarshal(message, &trade)
					if err != nil {
						i.Logger.Error("Unmarshal error:", err)
						continue
					}

					// price, err := strconv.ParseFloat(trade.Price, 64)
					// parse big float
					price, _, err := big.ParseFloat(trade.Price, 10, 64, big.ToNearestEven)
					if err != nil {
						i.Logger.Error(err.Error())
						continue
					}
					i.Mutex.Lock()
					i.SymbolPrices[trade.Symbol] = price
					i.Mutex.Unlock()

					i.ScanOpportunities(triade)
				}
			}
		}(triade)
	}
	wg.Wait()
}
