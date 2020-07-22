package main

import (
	"context"
	"github.com/rainu/launchpad-super-trigger/cmd/lst/config"
	"github.com/rainu/launchpad-super-trigger/gfx"
	launchpad "github.com/rainu/launchpad-super-trigger/pad"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {
	//initialise our global logger
	logger, _ := zap.NewDevelopment(
		zap.AddStacktrace(zap.FatalLevel), //disable stacktrace for level lower than fatal
	)
	zap.ReplaceGlobals(logger)
	defer zap.L().Sync()

	configContent := `
actors:
	rest:
		test:
			url: "http://localhost:1312"
layout:
	pages:
		0:
			trigger:
				"0,0":
					actor: test`
	configContent = strings.ReplaceAll(configContent, "\t", " ")

	dispatcher, err := config.ConfigureDispatcher(strings.NewReader(configContent))
	dispatcher.AddPageHandler(&launchpad.SimpleHandler{
		TriggerFn: func(lighter launchpad.Lighter, page launchpad.PageNumber, x int, y int) error {

			renderer := gfx.Renderer{lighter}
			renderer.Boom(x, y, launchpad.ColorHighGreen, 50*time.Millisecond)

			return nil
		},
		PageEnterFn: func(lighter launchpad.Lighter, page launchpad.PageNumber) error {
			renderer := gfx.Renderer{lighter}

			renderer.VerticalProgressbar(0, 75, gfx.AscDirection, launchpad.ColorHighGreen, launchpad.ColorNormalRed)

			return nil
		},
	}, 1, 3)

	if err != nil {
		zap.L().Fatal("error while opening setup launchpad configuration: %v", zap.Error(err))
	}

	pad, err := launchpad.NewLaunchpadSuperTrigger(dispatcher.Handle)
	if err != nil {
		zap.L().Fatal("error while opening connection to launchpad: %v", zap.Error(err))
	}
	defer pad.Close()

	//reacting to signals (interrupt)
	var signals = make(chan os.Signal, 1)
	defer close(signals)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	ctx, cancelFunc := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		pad.Run(ctx)
	}()

	//wait for interrupt
	<-signals

	cancelFunc()
	wg.Wait()
}
