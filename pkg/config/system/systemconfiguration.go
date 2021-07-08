package system

import (
	"os"

	"devcircus.com/cerberus/pkg/util/logger"

	"github.com/jinzhu/configor"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// SystemConfiguration configuration
var SystemConfiguration = struct {
	AppName  string `yaml:"appname"`
	Version  string `yaml:"version"`
	LogLevel string `yaml:"loglevel"`
}{}

func init() {
	// Load configuration file
	configor.Load(&SystemConfiguration, "config.yml")
	// Configure logger
	logger.ConfigLogger(SystemConfiguration.LogLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Debug().
		Interface("SystemConfiguration", SystemConfiguration).
		Msg("System Configuration:")
}

// GetAppName Returns the app name value
func GetAppName() string {
	return SystemConfiguration.AppName
}

// SetAppName set the app name value
func SetAppName(appName string) {
	SystemConfiguration.AppName = appName
}

// GetVersion returns the app version
func GetVersion() string {
	return SystemConfiguration.Version
}

// SetVersion set the app version
func SetVersion(version string) {
	SystemConfiguration.Version = version
}

// GetLogLevel returns the log level value
func GetLogLevel() string {
	return SystemConfiguration.LogLevel
}
