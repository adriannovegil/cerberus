package execute

import (
	"flag"
	"os"
	"syscall"

	"github.com/rs/zerolog/log"
	daemon "github.com/sevlyar/go-daemon"
)

var (
	// Stop signal
	Stop = make(chan struct{})
	// Done signal
	Done = make(chan struct{})
	// Reload signal
	Reload      = make(chan struct{})
	signalFlags = flag.String("s", "", `Send signal to the daemon:
quit — graceful shutdown
stop — fast shutdown
reload — reloading the configuration file`)
)

func termHandler(sig os.Signal) error {
	log.Debug().Msg("Termination requested")
	if sig == syscall.SIGQUIT {
		Done <- struct{}{}
	}  else {
		Stop <- struct{}{}
	}
	return daemon.ErrStop
}

func reloadHandler(sig os.Signal) error {
	log.Debug().Msg("Configuration reload requested")
	Reload <- struct{}{}
	return nil
}
