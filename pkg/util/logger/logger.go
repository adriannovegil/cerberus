package logger

import (
	"github.com/rs/zerolog"
)

const (
	//LevelPanic constant
	LevelPanic = 5
	//LevelFatal constant
	LevelFatal = 4
	//LevelError constant
	LevelError = 3
	//LevelWarn constant
	LevelWarn = 2
	//LevelInfo constant
	LevelInfo = 1
	//LevelDebug constant
	LevelDebug = 0
	//LevelTrace constant
	LevelTrace = -1
	//LevelNoLevel constant
	LevelNoLevel = -2
)

// ConfigLogger Configure the system logger
func ConfigLogger(logLevel string) {
	// Configure logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	// Change the default log level
	switch toLogLevelValue(logLevel) {
	case LevelPanic:
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case LevelFatal:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case LevelError:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case LevelWarn:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case LevelInfo:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case LevelDebug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case LevelTrace:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
}

// Convert the log leves to the string value
func toLogLevelString(logLevel int) string {
	switch logLevel {
	case LevelPanic:
		return "panic"
	case LevelFatal:
		return "fatal"
	case LevelError:
		return "error"
	case LevelWarn:
		return "warn"
	case LevelInfo:
		return "info"
	case LevelDebug:
		return "debug"
	case LevelTrace:
		return "trace"
	case LevelNoLevel:
		return ""
	}
	return ""
}

// Convert the log leves to the int value
func toLogLevelValue(logLevel string) int {
	if logLevel == "panic" {
		return LevelPanic
	} else if logLevel == "fatal" {
		return LevelFatal
	} else if logLevel == "error" {
		return LevelError
	} else if logLevel == "warn" {
		return LevelWarn
	} else if logLevel == "info" {
		return LevelInfo
	} else if logLevel == "debug" {
		return LevelDebug
	} else if logLevel == "trace" {
		return LevelTrace
	} else {
		return LevelNoLevel
	}
}
