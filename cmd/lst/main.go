package main

import (
	"context"
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
	//initialise our global logger
	logger, _ := zap.NewDevelopment(
		zap.AddStacktrace(zap.FatalLevel), //disable stacktrace for level lower than fatal
	)
	zap.ReplaceGlobals(logger)
	defer zap.L().Sync()

	dispatcher := launchpad.TriggerDispatcher{}
	dispatcher.AddPageHandler(&launchpad.SimpleHandler{
		TriggerFn: func(lighter launchpad.Lighter, page launchpad.PageNumber, x int, y int) error {

			renderer := gfx.Renderer{lighter}
			renderer.Star(x, y, launchpad.ColorHighGreen, 250*time.Millisecond)

			return nil
		},
		PageEnterFn: func(lighter launchpad.Lighter, page launchpad.PageNumber) error {
			renderer := gfx.Renderer{lighter}

			renderer.Boom(4, 4, launchpad.ColorHighGreen, 50*time.Millisecond)

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
