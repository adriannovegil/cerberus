package main

import (
	"os"
	"runtime"

	"devcircus.com/cerberus/pkg/cli"
	"devcircus.com/cerberus/pkg/config/system"

	"github.com/rs/zerolog/log"
)

// Main function
func main() {
	log.Debug().Msgf("datanerd version: %s", system.GetVersion())
	log.Debug().Msgf("go runtime version: %s", runtime.Version())
	log.Debug().Msgf("CLI args: %#v", os.Args)
	cli.Execute()
}
