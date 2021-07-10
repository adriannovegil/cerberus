package execute

import (
	"flag"
	"os"
	"syscall"

	"github.com/rs/zerolog/log"
	daemon "github.com/sevlyar/go-daemon"
)

var (
	stop        = make(chan struct{})
	done        = make(chan struct{})
	signalFlags = flag.String("s", "", `Send signal to the daemon:
quit — graceful shutdown
stop — fast shutdown
reload — reloading the configuration file`)
)

func termHandler(sig os.Signal) error {
	log.Info().Msg("terminating...")
	stop <- struct{}{}
	if sig == syscall.SIGQUIT {
		<-done
	}
	return daemon.ErrStop
}

func reloadHandler(sig os.Signal) error {
	log.Info().Msg("configuration reloaded")
	return nil
}
