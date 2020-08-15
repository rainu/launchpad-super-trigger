package actor

import (
	"github.com/rainu/launchpad-super-trigger/actor"
	"github.com/rainu/launchpad-super-trigger/config"
	"github.com/rainu/launchpad-super-trigger/pad"
	"time"
)

func buildGfxBlink(actors map[string]actor.Actor, gfxActors map[string]config.GfxBlinkActor) {
	for actorName, gActor := range gfxActors {
		handler := &actor.GfxBlink{
			ColorOn:  colorOrPanic(gActor.ColorOn),
			ColorOff: colorOrDefault(gActor.ColorOff, pad.ColorOff),
			Interval: gActor.Interval,
			Duration: gActor.Duration,
		}
		if handler.Interval == 0 {
			handler.Interval = 1 * time.Second
		}

		actors[actorName] = handler
	}
}

func colorOrDefault(color config.Color, defaultColor pad.Color) pad.Color {
	c, err := color.Color()
	if err != nil {
		return defaultColor
	}
	return c
}

func colorOrPanic(color config.Color) pad.Color {
	c, err := color.Color()
	if err != nil {
		panic(err)
	}
	return c
}
