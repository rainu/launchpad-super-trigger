package main

import (
	"context"
	"fmt"
	"github.com/rainu/launchpad-super-trigger/cmd/lst/config"
	launchpad "github.com/rainu/launchpad-super-trigger/pad"
	driver "gitlab.com/gomidi/rtmididrv"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

	dispatcher, sensors, generalSettings, err := config.ConfigureDispatcher(configFile)
	if err != nil {
		zap.L().Fatal("error while opening setup launchpad configuration: %v", zap.Error(err))
	}

	d, err := driver.New()
	if err != nil {
		zap.L().Fatal("Unable to load midi driver: %s", zap.Error(err))
	}

	pad, err := launchpad.NewLaunchpadSuperTrigger(d, dispatcher.Handle)
	if err != nil {
		zap.L().Fatal("error while opening connection to launchpad: %v", zap.Error(err))
	}
	defer pad.Close()

	err = pad.Initialise(
		generalSettings.StartPage.AsInt(),
		generalSettings.NavigationMode,
	)
	if err != nil {
		zap.L().Fatal("error while initialize launchpad: %v", zap.Error(err))
	}

	//reacting to signals (interrupt)
	var signals = make(chan os.Signal, 1)
	var connectionLost = make(chan bool, 1)
	defer close(signals)
	defer close(connectionLost)
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

	go func() {
		pad.WaitForConnectionLost(ctx)
		connectionLost <- true
	}()

	//wait for interrupt or connection lost
	select {
	case <-signals:
		zap.L().Info("Interrupt signal received.")
	case <-connectionLost:
		zap.L().Info("Connection to Launchpad lost.")
	}

	cancelFunc()
	wg.Wait()
}
