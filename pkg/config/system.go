package config

import (
	"os"

	"github.com/jinzhu/configor"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"devcircus.com/cerberus/pkg/target/request"
	"devcircus.com/cerberus/pkg/util/logger"
)

// System configuration
var System = struct {
	AppName     string  `yaml:"appname"`
	Version     string  `yaml:"version"`
	LogLevel    string  `yaml:"loglevel"`
	Concurrency int     `yaml:"concurrency"`
	Targets     Targets `yaml:"targets"`
}{}

// Targets to control
type Targets struct {
	Requests []request.Config `yaml:"requests"`
}

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

	if err := validate(System.Targets.Requests); err != nil {
		log.Fatal().Err(err).Msg("Invalid Request data in config file")
	}
	System.Targets.Requests = generateAndAssignIdsForRequests(System.Targets.Requests)
}

// ReLoad system configuration
func ReLoad() {
	load()
}
