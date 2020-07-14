package main

import (
	"context"
	launchpad "github.com/rainu/launchpad-super-trigger/pad"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	//initialise our global logger
	logger, _ := zap.NewDevelopment(
		zap.AddStacktrace(zap.FatalLevel), //disable stacktrace for level lower than fatal
	)
	zap.ReplaceGlobals(logger)
	defer zap.L().Sync()

	pad, err := launchpad.NewLaunchpadSuperTrigger()
	if err != nil {
		zap.L().Fatal("error while openning connection to launchpad: %v", zap.Error(err))
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
