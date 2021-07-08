package config

import (
	"os"

	"devcircus.com/cerberus/pkg/util/logger"

	"github.com/jinzhu/configor"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// System configuration
var System = struct {
	AppName  string `yaml:"appname"`
	Version  string `yaml:"version"`
	LogLevel string `yaml:"loglevel"`
	Port     int    `yaml:"port"`
}{}

func init() {
	// Load configuration file
	configor.Load(&System, "config.yml")
	// Configure logger
	logger.ConfigLogger(System.LogLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Debug().
		Interface("System", System).
		Msg("System Configuration:")
}
