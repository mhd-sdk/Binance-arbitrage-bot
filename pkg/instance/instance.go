package instance

import (
	"DeltA/pkg/binance"
	"DeltA/pkg/helpers"
	"DeltA/pkg/models"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"strconv"

	"github.com/gorilla/websocket"
)

type Instance struct {
	ExchangeInfo        models.ExchangeInfo
	SymbolPrices        map[string]float64
	StartingStableAsset string
	Triades             []models.Triade
	Logger              *slog.Logger
}

func NewInstance(logger *slog.Logger) *Instance {
	return &Instance{Logger: logger}
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
	// for now, we will only watch the first triade
	triade := i.Triades[5]

	pairSymbols := []string{}
	for _, symbol := range triade.Symbols {
		pairSymbols = append(pairSymbols, symbol.Symbol)
	}

	url := helpers.BuildStreamURL(pairSymbols)

	i.Logger.Info(fmt.Sprintf("Connecting to %s", url))

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal(err, url)
	}

	defer conn.Close()
	i.Logger.Info("Watching for arbitrage opportunities... (" + triade.Assets[0] + "/" + triade.Assets[1] + "/" + triade.Assets[2] + ")")

	i.ScanOpportunities(triade)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}

		var trade models.Trade
		err = json.Unmarshal(message, &trade)
		if err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}
		price, err := strconv.ParseFloat(trade.Price, 64)
		if err != nil {
			i.Logger.Error(err.Error())
			continue
		}

		i.SymbolPrices[trade.Symbol] = price

		// if i.SymbolPrices[triade.Symbols[0].Symbol] == 0 || i.SymbolPrices[triade.Symbols[1].Symbol] == 0 || i.SymbolPrices[triade.Symbols[2].Symbol] == 0 {
		// 	continue
		// }

		i.ScanOpportunities(triade)
	}
}
