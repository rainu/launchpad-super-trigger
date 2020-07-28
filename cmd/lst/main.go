package main

import (
	"context"
	"fmt"
	"github.com/rainu/launchpad-super-trigger/cmd/lst/config"
	"github.com/rainu/launchpad-super-trigger/gfx"
	launchpad "github.com/rainu/launchpad-super-trigger/pad"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	LoadArgs()

	//initialise our global logger
	logger, _ := zap.NewDevelopment(
		zap.AddStacktrace(zap.FatalLevel), //disable stacktrace for level lower than fatal
	)
	zap.ReplaceGlobals(logger)
	defer zap.L().Sync()

	configFile, err := os.Open(*Args.ConfigFile)
	if err != nil {
		zap.L().Fatal("error while read configuration: %v", zap.Error(err))
	}

	dispatcher, sensors, err := config.ConfigureDispatcher(configFile)
	if err != nil {
		zap.L().Fatal("error while opening setup launchpad configuration: %v", zap.Error(err))
	}

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

	//start sensors
	for sensorName, curSensor := range sensors {
		zap.L().Info(fmt.Sprintf("Start sensor: %s", sensorName))
		go curSensor.Run(ctx)
	}

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
