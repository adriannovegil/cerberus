package main

import (
	"os"
	"runtime"

	"devcircus.com/cerberus/pkg/cli"
	"devcircus.com/cerberus/pkg/config"

	"github.com/rs/zerolog/log"
)

// Main function
func main() {

	log.Debug().Msgf("cerberus version: %s", config.Config.Version)
	log.Debug().Msgf("go runtime version: %s", runtime.Version())
	log.Debug().Msgf("CLI args: %#v", os.Args)
	cli.Execute()
}
