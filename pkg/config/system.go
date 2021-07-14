package config

import (
	"os"

	"github.com/jinzhu/configor"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"devcircus.com/cerberus/pkg/requests"
	"devcircus.com/cerberus/pkg/util/logger"
)

// System configuration
var System = struct {
	AppName     string                   `yaml:"appname"`
	Version     string                   `yaml:"version"`
	LogLevel    string                   `yaml:"loglevel"`
	Concurrency int                      `yaml:"concurrency"`
	Requests    []requests.RequestConfig `yaml:"requests"`
}{}

func init() {
	load()
	// Configure logger
	logger.ConfigLogger(System.LogLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Debug().
		Interface("System", System).
		Msg("System Configuration:")
}

// load system configuration
func load() {
	configor.Load(&System, "config.yml")
	System.Requests = validateAndCreateIdsForRequests(System.Requests)
}

// ReLoad system configuration
func ReLoad() {
	load()
}
