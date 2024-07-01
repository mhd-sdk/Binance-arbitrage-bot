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
	"math"
	"math/big"
	"os"
	"strconv"
	"sync"
	"time"

	binance_connector "github.com/binance/binance-connector-go"
	"github.com/gorilla/websocket"
	"github.com/kr/pretty"
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

	// _, err = i.BinanceConn.NewCreateOrderService().Symbol("DGBBTC").
	// 	Side("BUY").Type("MARKET").Quantity(2885).
	// 	Do(context.Background())
	// if err != nil {
	// 	i.Logger.Slog.Error(err.Error())
	// }
	// os.Exit(0)

	i.Logger.Slog.Info("Starting DeltΔ  BOT")

	i.StartingStableAsset = "BTC"

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
	// wg.Add(1)

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

func ComputeTriangularOrders(i *Instance, orderPairs OrderPairs, gainPercentage *big.Float) {
	if orderCount == 1 || isOrdering {
		return
	}
	isOrdering = true
	orderCount++

	logs := "Triangular arbitrage opportunity found : " + orderPairs[0].From + "/" + orderPairs[1].From + "/" + orderPairs[2].From + " Gain: " + gainPercentage.Text('f', -1) + "%"
	i.Logger.Slog.Info(logs)

	var Balance = Balance{}
	logBalance := map[string]float64{}

	account, err := i.BinanceConn.NewGetAccountService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, balance := range account.Balances {
		if balance.Asset == i.StartingStableAsset {
			float, _ := strconv.ParseFloat(balance.Free, 64)
			Balance[i.StartingStableAsset] = big.NewFloat(float)
			logBalance[i.StartingStableAsset] = float
			if err != nil {
				i.Logger.Slog.Error(err.Error())
				return
			}
			break
		}
	}

	for _, pair := range orderPairs {
		symbolPrice, _ := pair.SymbolPrice.Float64()
		if helpers.CheckAssetIsBase(pair.To, pair.Symbol) {
			Balance[pair.To] = new(big.Float).Quo(Balance[pair.From], pair.SymbolPrice)
			logBalance[pair.To] = logBalance[pair.From] / symbolPrice
		} else {
			Balance[pair.To] = new(big.Float).Mul(Balance[pair.From], pair.SymbolPrice)
			logBalance[pair.To] = logBalance[pair.From] * symbolPrice
		}
	}

	pretty.Println(logBalance)

	// os.Exit(0)

	i.Logger.Slog.Info(logs)

	for _, orderPair := range orderPairs {
		// Step 1: Get the account balance
		// account, err := i.BinanceConn.NewGetAccountService().Do(context.Background())
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }

		// var fromBalance float64
		// for _, balance := range account.Balances {
		// 	if balance.Asset == orderPair.From {
		// 		fromBalance, err = strconv.ParseFloat(balance.Free, 64)
		// 		if err != nil {
		// 			i.Logger.Slog.Error(err.Error())
		// 			return
		// 		}
		// 		break
		// 	}
		// }

		// safetyMargin := 0.5
		// fromBalance = fromBalance * safetyMargin

		// quantity, _ := new(big.Float).Quo(big.NewFloat(fromBalance), price).SetPrec(uint(math.Ceil(float64(6) * 3.32))).Float64()
		if helpers.CheckAssetIsBase(orderPair.To, orderPair.Symbol) {
			// check symbolFilters
			// assetBuyedQty, _ := big.NewFloat(0).Quo(big.NewFloat(fromBalance), orderPair.SymbolPrice).Float64()
			// fmt.Println(assetBuyedQty)
			// for _, filter := range orderPair.Symbol.Filters {
			// 	if filter.FilterType == "LOT_SIZE" {
			// 		// changer fromBalance pour qu'il soit un multiple de stepSize
			// 		stepSize, _ := strconv.ParseFloat(filter.StepSize, 64)
			// 		assetBuyedQty = math.Floor(assetBuyedQty/stepSize) * stepSize
			// 	}
			// }
			// fmt.Println(assetBuyedQty)
			// fmt.Println(assetBuyedQty)

			// quoteQty, _ := big.NewFloat(0).Quo(orderPair.SymbolPrice, big.NewFloat(fromBalance)).Float64()

			buyingQty, _ := Balance[orderPair.To].Float64()
			for _, filter := range orderPair.Symbol.Filters {
				if filter.FilterType == "LOT_SIZE" {
					stepSize, _ := strconv.ParseFloat(filter.StepSize, 64)
					buyingQty = math.Floor(buyingQty/stepSize) * stepSize
				}
			}

			buyingQty = helpers.FloorToPrecision(buyingQty, int(orderPair.Symbol.BaseAssetPrecision))
			i.Logger.Slog.Info(fmt.Sprintf("Order BUY Symbol:%s price:%s qty:%.10f", orderPair.Symbol.Symbol, orderPair.SymbolPrice.String(), buyingQty))
			_, err = i.BinanceConn.NewCreateOrderService().Symbol(orderPair.Symbol.Symbol).
				Side("BUY").Type("MARKET").Quantity(buyingQty).
				Do(context.Background())
			if err != nil {
				// fmt.Println(binance_connector.PrettyPrint(order))
				i.Logger.Slog.Error(err.Error())
			}

		} else {

			// fromBalance = helpers.FloorToPrecision(fromBalance, int(orderPair.Symbol.BaseAssetPrecision))
			sellingQty, _ := Balance[orderPair.From].Float64()

			for _, filter := range orderPair.Symbol.Filters {
				if filter.FilterType == "LOT_SIZE" {
					stepSize, _ := strconv.ParseFloat(filter.StepSize, 64)
					sellingQty = math.Floor(sellingQty/stepSize) * stepSize
				}
			}

			sellingQty = helpers.FloorToPrecision(sellingQty, int(orderPair.Symbol.BaseAssetPrecision))

			i.Logger.Slog.Info(fmt.Sprintf("Order SELL Symbol:%s price:%s qty:%.10f", orderPair.Symbol.Symbol, orderPair.SymbolPrice.String(), sellingQty))

			_, err = i.BinanceConn.NewCreateOrderService().Symbol(orderPair.Symbol.Symbol).
				Side("SELL").Type("MARKET").Quantity(sellingQty).
				Do(context.Background())
			if err != nil {
				i.Logger.Slog.Error(err.Error())
			}
		}

		// time.Sleep(1 * time.Second)
	}
	os.Exit(0)
	isOrdering = false
}
