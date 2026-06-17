package config

import (
	"github.com/patharetush/incident-pilot/shared/config"
)

type Config = config.Config

const (
	TransportHTTP  = config.TransportHTTP
	TransportStdio = config.TransportStdio
)

var defaults = config.Defaults{
	ServerName:    "incident-pilot-monitoring",
	ServerVersion: "0.1.0",
	HTTPAddr:      ":8081",
	LogFilename:   "/tmp/monitoring-mcp.log",
}

func Default() *Config {
	return config.Default(defaults)
}

func Load() *Config {
	return config.Load(defaults)
}
