package main

import (
	"prompter/internal/domain/entity"
	"prompter/internal/domain/port"
	"prompter/internal/pkg/registry"
)

func main() {
	messengersRegistry := registry.Registry[entity.Provider, port.Messenger](make(map[entity.Provider]port.Messenger))

	var d port.MessengerRegistry = messengersRegistry
	_, _ = d.Get(entity.ProviderHH)
}
