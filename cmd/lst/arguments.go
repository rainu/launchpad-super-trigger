package main

import (
	"flag"
	"go.uber.org/zap"
)

type applicationArgs struct {
	ConfigPath *string
	DebugPort  *int
}

var Args applicationArgs

func LoadArgs() {
	Args = applicationArgs{
		ConfigPath: flag.String("config", "./config.yaml", "The configuration file/directory"),
		DebugPort:  flag.Int("debug", -1, "The port for the (pprof) debug endpoint"),
	}
	flag.Parse()

	if *Args.ConfigPath == "" {
		zap.L().Fatal("Topic configuration file/directory is missing!")
	}
}
