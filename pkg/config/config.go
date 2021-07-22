package config

import (
	"os"

	"github.com/jinzhu/configor"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"devcircus.com/cerberus/pkg/fallback"
	"devcircus.com/cerberus/pkg/metrics/prometheus"
	"devcircus.com/cerberus/pkg/target/request"
	"devcircus.com/cerberus/pkg/util/logger"
)

// Config configuration
var Config = struct {
	AppName     string `yaml:"appname"`
	Version     string `yaml:"version"`
	LogLevel    string `yaml:"loglevel"`
	Concurrency int    `yaml:"concurrency"`
	Metrics     Metrics
	Fallbacks   []fallback.Config `yaml:"fallbacks"`
	Targets     Targets
}{}

// Targets to control
type Targets struct {
	Requests []request.Config `yaml:"requests"`
}

// Metrics recorders
type Metrics struct {
	Prometheus prometheus.Config `yaml:"prometheus"`
}

func init() {
	load()
	// Configure logger
	logger.ConfigLogger(Config.LogLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Debug().
		Interface("Config", Config).
		Msg("System Configuration:")
}

// load system configuration
func load() {
	configor.Load(&Config, "config.yml")

	if err := validate(Config.Targets.Requests); err != nil {
		log.Fatal().Err(err).Msg("Invalid Request data in config file")
	}
	Config.Targets.Requests = generateAndAssignIdsForRequests(Config.Targets.Requests)
}

// ReLoad system configuration
func ReLoad() {
	load()
}

// GetFallbackCOnfigurationByName retutrn the fallback configuration by name
func GetFallbackCOnfigurationByName(fallbackActionName string) *fallback.Config {
	for _, fallbackConfig := range Config.Fallbacks {
		if fallbackConfig.Name == fallbackActionName {
			return &fallbackConfig
		}
	}
	return nil
}
