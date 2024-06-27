package main

import (
	"DeltA/pkg/instance"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, nil))
	logger.Info("Starting DeltΔ  BOT")

	instance := instance.NewInstance(logger)
	err := instance.Init()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	instance.Watch()
}
