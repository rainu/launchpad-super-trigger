package main

import (
	"flag"
	"go.uber.org/zap"
)

type applicationArgs struct {
	ConfigPath *string
}

var Args applicationArgs

func LoadArgs() {
	Args = applicationArgs{
		ConfigPath: flag.String("config", "./config.yaml", "The configuration file/directory"),
	}
	flag.Parse()

	if *Args.ConfigPath == "" {
		zap.L().Fatal("Topic configuration file/directory is missing!")
	}
}
