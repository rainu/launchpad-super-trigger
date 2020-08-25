package main

import (
	"context"
	"fmt"
	"github.com/rainu/launchpad-super-trigger/cmd/lst/config"
	triggerConf "github.com/rainu/launchpad-super-trigger/config"
	launchpad "github.com/rainu/launchpad-super-trigger/pad"
	"github.com/rainu/launchpad-super-trigger/sensor"
	driver "gitlab.com/gomidi/rtmididrv"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
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

	//start pprof if needed
	if *Args.DebugPort > 0 {
		zap.L().Info(fmt.Sprintf("Start pprof debug endpoint :%d", *Args.DebugPort))
		go func() {
			log.Println(http.ListenAndServe(fmt.Sprintf(":%d", *Args.DebugPort), nil))
		}()
	}

	dispatcher, sensors, generalSettings := initialiseConfig()

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

func initialiseConfig() (*launchpad.TriggerDispatcher, map[string]sensor.Sensor, triggerConf.General) {
	info, err := os.Stat(*Args.ConfigPath)
	if os.IsNotExist(err) {
		zap.L().Fatal("Config file or directory not found!")
	}

	var files []*os.File

	if info.IsDir() {
		dir, err := ioutil.ReadDir(*Args.ConfigPath)
		if err != nil {
			zap.L().Fatal("could not read config directory: %v", zap.Error(err))
		}
		files = make([]*os.File, 0, len(dir))

		for _, fileInfo := range dir {
			if !fileInfo.IsDir() {
				configFile, err := os.Open(fmt.Sprintf("%s/%s", *Args.ConfigPath, fileInfo.Name()))
				if err != nil {
					zap.L().Fatal("error while read configuration: %v", zap.Error(err))
				}

				files = append(files, configFile)
			}
		}
	} else {
		configFile, err := os.Open(*Args.ConfigPath)
		if err != nil {
			zap.L().Fatal("error while read configuration: %v", zap.Error(err))
		}

		files = []*os.File{configFile}
	}
	defer func() {
		//close all config files
		for _, file := range files {
			file.Close()
		}
	}()

	reader := make([]io.Reader, 0, len(files))
	for _, file := range files {
		reader = append(reader, file)
	}

	dispatcher, sensors, generalSettings, err := config.ConfigureDispatcher(reader...)
	if err != nil {
		zap.L().Fatal("error while opening setup launchpad configuration: %v", zap.Error(err))
	}

	return dispatcher, sensors, generalSettings
}
