package execute

import (
	"syscall"

	"github.com/rs/zerolog/log"
	daemon "github.com/sevlyar/go-daemon"
)

// Daemon execution mode
func Daemon(proc func()) {
	daemon.AddCommand(daemon.StringFlag(signalFlags, "quit"), syscall.SIGQUIT, termHandler)
	daemon.AddCommand(daemon.StringFlag(signalFlags, "stop"), syscall.SIGTERM, termHandler)
	daemon.AddCommand(daemon.StringFlag(signalFlags, "stop"), syscall.SIGKILL, termHandler)
	daemon.AddCommand(daemon.StringFlag(signalFlags, "reload"), syscall.SIGHUP, reloadHandler)

	ctx := &daemon.Context{
		PidFileName: "cerberus.pid",
		PidFilePerm: 0644,
		WorkDir:     "./",
		Umask:       027,
	}

	child, err := ctx.Reborn()
	if err != nil {
		log.Fatal().Msgf("An error occured while trying to reborn daemon %s", err.Error())
	}
	if child != nil {
		return
	}
	defer ctx.Release()

	go proc()

	err = daemon.ServeSignals()
	if err != nil {
		log.Error().Msgf("Error %s", err.Error())
	}

	log.Info().Msg("daemon terminated")
}
