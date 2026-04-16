package main

import (
	"hr-bot-ddd-example/internal/domain/entity"
	"hr-bot-ddd-example/internal/domain/port"
	"hr-bot-ddd-example/internal/pkg/registry"
)

func main() {
	messengersRegistry := registry.Registry[entity.Provider, port.Messenger](make(map[entity.Provider]port.Messenger))

	var d port.MessengerRegistry = messengersRegistry
	_, _ = d.Get(entity.ProviderHH)
}
