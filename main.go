package main

import (
	"DeltA/pkg/instance"
	"log/slog"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	instance := instance.NewInstance()
	err = instance.Init()
	if err != nil {
		instance.Logger.Slog.Error(err.Error())
		return
	}
	instance.Watch()
}
