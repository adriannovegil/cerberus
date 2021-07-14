package server

import (
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"devcircus.com/cerberus/pkg/execute"
	"devcircus.com/cerberus/pkg/worker"
)

const defaultTick = 60 * time.Second

var (
	supervisor *worker.Supervisor
)

// NewCmdServer start the server
func NewCmdServer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "start the system server",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info().Msgf("Starting daemon with pid: %d", os.Getpid())
			isDaemon, _ := cmd.Flags().GetBool("daemon")
			supervisor = worker.NewSupervisor()
			run(os.Getpid(), isDaemon)
		},
	}
	addFlags(cmd)
	return cmd
}

// Add flags to the command
func addFlags(cmd *cobra.Command) {
	cmd.Flags().BoolP("daemon", "d", false, "Daemon execution")
}

func run(pidFile int, isDaemon bool) {
	if isDaemon {
		log.Info().Msg("Executing in daemon mode ...")
		execute.Daemon(supervisor.Run)
	} else {
		log.Info().Msg("Executing in interactive mode ...")
		execute.Interactive(supervisor.Run)
	}
}
