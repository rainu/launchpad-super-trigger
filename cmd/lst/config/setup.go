package config

import (
	configActor "github.com/rainu/launchpad-super-trigger/cmd/lst/config/actor"
	"github.com/rainu/launchpad-super-trigger/cmd/lst/config/connection"
	"github.com/rainu/launchpad-super-trigger/config"
	"github.com/rainu/launchpad-super-trigger/pad"
	"io"
)

func ConfigureDispatcher(configReader io.Reader) (*pad.TriggerDispatcher, error) {
	parsedConfig, err := config.ReadConfig(configReader)
	if err != nil {
		return nil, err
	}

	dispatcher := &pad.TriggerDispatcher{}
	connections := connection.BuildMqttConnection(parsedConfig)
	actors := configActor.BuildActors(parsedConfig, connections)

	for pageNumber, page := range parsedConfig.Layout.Pages {
		handler := &pageHandler{
			pageNumber: pad.PageNumber(pageNumber),
			page:       page,
		}
		handler.Init(actors)
		dispatcher.AddPageHandler(handler, handler.pageNumber)
	}

	return dispatcher, nil
}
