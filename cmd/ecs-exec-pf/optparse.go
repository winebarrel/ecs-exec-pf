package main

import (
	"log"
	"runtime/debug"

	"github.com/integrii/flaggy"
)

const (
	Description = "Port forwarding using the ECS task container. (aws-cli wrapper)"
)

type Options struct {
	Cluster   string
	Task      string
	Container string
	Port      uint16
	LocalPort uint16
	Debug     bool
}

func getVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		return info.Main.Version
	}

	return "unknown"
}

func parseArgs() *Options {
	opts := &Options{}

	flaggy.SetDescription(Description)
	flaggy.SetVersion(getVersion())
	flaggy.String(&opts.Cluster, "c", "cluster", "ECS cluster name.")
	flaggy.String(&opts.Task, "t", "task", "ECS task ID.")
	flaggy.String(&opts.Container, "n", "container", "Container name in ECS task.")
	flaggy.UInt16(&opts.Port, "p", "port", "Target remote port.")
	flaggy.UInt16(&opts.LocalPort, "l", "local-port", "Client local port.")
	flaggy.Parse()

	if opts.Cluster == "" {
		log.Fatal("'--cluster' is required")
	}

	if opts.Task == "" {
		log.Fatal("'--task' is required")
	}

	if opts.Port == 0 {
		log.Fatal("'--port' is required")
	}

	if opts.LocalPort == 0 {
		log.Fatal("'--local-port' is required")
	}

	return opts
}
