package instance

import (
	"DeltA/pkg/binance"
	"DeltA/pkg/helpers"
	"DeltA/pkg/models"
	"fmt"
	"log/slog"

	"github.com/kr/pretty"
)

type Instance struct {
	ExchangeInfo binance.ExchangeInfo
	Triades      []models.Triade
	Logger       *slog.Logger
}

func NewInstance(logger *slog.Logger) *Instance {
	return &Instance{Logger: logger}
}

func (i *Instance) Init() (err error) {
	i.Logger.Info("Loading exchange info from Binance...")
	i.ExchangeInfo, err = binance.GetExchangeInfo()
	if err != nil {
		return err
	}
	i.Logger.Info(fmt.Sprintf("Number of pairs: %d", len(i.ExchangeInfo)))
	triades := helpers.BuildTriades(i.ExchangeInfo)
	i.Logger.Info(fmt.Sprintf("Number of triades: %d", len(triades)))
	pretty.Println(triades[10])
	return nil
}
