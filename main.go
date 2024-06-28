package main

import (
	"DeltA/pkg/instance"
)

func main() {

	instance := instance.NewInstance()
	err := instance.Init()
	if err != nil {
		instance.Logger.Slog.Error(err.Error())
		return
	}
	instance.Watch()
}
