package instance

import (
	"DeltA/pkg/binance"
	"DeltA/pkg/helpers"
	"DeltA/pkg/logging"
	"DeltA/pkg/models"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"strconv"
	"sync"
	"time"

	binance_connector "github.com/binance/binance-connector-go"
	"github.com/gorilla/websocket"
)

type Instance struct {
	BinanceConn  *binance_connector.Client
	Mutex        *sync.Mutex
	ExchangeInfo *binance_connector.ExchangeInfoResponse
	// SymbolPrices        map[string]*big.Float
	SymbolPrices        map[string]*big.Float
	StartingStableAsset string
	Triades             []models.Triade
	Logger              *logging.Logger
}

func NewInstance() *Instance {
	return &Instance{}
}

func (i *Instance) Init() (err error) {
	BINANCE_API_KEY := os.Getenv("BINANCE_API_KEY")
	BINANCE_SECRET_KEY := os.Getenv("BINANCE_SECRET_KEY")
	i.BinanceConn = binance_connector.NewClient(BINANCE_API_KEY, BINANCE_SECRET_KEY, "https://api.binance.com")

	i.Mutex = new(sync.Mutex)

	logger, err := logging.NewLogger()
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	i.Logger = logger

	i.Logger.Slog.Info("Starting DeltΔ  BOT")

	i.StartingStableAsset = "USDC"

	i.Logger.Slog.Info("Loading exchange info from Binance...")
	// i.ExchangeInfo, err = binance.GetExchangeInfo()
	i.ExchangeInfo, err = i.BinanceConn.NewExchangeInfoService().Do(context.Background())

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

			for {
				conn, _, err := websocket.DefaultDialer.Dial(url, nil)
				if err != nil {
					i.Logger.Slog.Error(fmt.Sprintf("Error connecting to %s: %s", url, err))
					time.Sleep(5 * time.Second) // Retry after 5 seconds
					continue
				}

				defer conn.Close()

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

					i.ScanOpportunities(triade)
				}
			}
		}(triade)
	}
	wg.Wait()
}

var orderCount = 0
var isOrdering = false

func (i *Instance) ComputeTriangularOrders(orderPairs OrderPairs) {
	if orderCount == 1 || isOrdering {
		return
	}
	isOrdering = true
	orderCount++
	i.Logger.Slog.Info("Computing triangular orders...")

	for _, orderPair := range orderPairs {
		// Step 1: Get the account balance
		account, err := i.BinanceConn.NewGetAccountService().Do(context.Background())
		if err != nil {
			fmt.Println(err)
			return
		}

		var fromBalance float64
		for _, balance := range account.Balances {
			if balance.Asset == orderPair.From {
				fromBalance, err = strconv.ParseFloat(balance.Free, 64)
				if err != nil {
					i.Logger.Slog.Error(err.Error())
					return
				}
				break
			}
		}

		price := orderPair.Price

		side := "BUY"
		if orderPair.To+orderPair.From != orderPair.Symbol.Symbol {
			side = "SELL"
			price = new(big.Float).Quo(big.NewFloat(1), price)
		}

		i.Logger.Slog.Info(fmt.Sprintf("Order %s from:%s to:%s price:%s quoteQty:%f symbol:%s", side, orderPair.From, orderPair.To, orderPair.Price.String(), fromBalance, orderPair.Symbol.Symbol))
		// order, err := i.BinanceConn.NewCreateOrderService().Symbol(orderPair.Symbol.Symbol).
		// 	Side(side).Type("MARKET").QuoteOrderQty(fromBalance).
		// 	Do(context.Background())
		// if err != nil {
		// 	binance_connector.PrettyPrint(order)
		// 	i.Logger.Slog.Error(err.Error())
		// 	return
		// }

	}
	isOrdering = false
}
