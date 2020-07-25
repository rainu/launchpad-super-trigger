package main

import (
	"flag"
	"go.uber.org/zap"
)

type applicationArgs struct {
	ConfigFile *string
}

var Args applicationArgs

func LoadArgs() {
	Args = applicationArgs{
		ConfigFile: flag.String("config", "./config.yaml", "The configuration file"),
	}
	flag.Parse()

	if *Args.ConfigFile == "" {
		zap.L().Fatal("Topic configuration file is missing!")
	}
}
