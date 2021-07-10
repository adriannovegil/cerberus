package server

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"devcircus.com/cerberus/pkg/execute"
)

const defaultTick = 60 * time.Second

// NewCmdServer start the server
func NewCmdServer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "start the system server",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info().Msg(fmt.Sprintf("Starting server with pid: %d", os.Getpid()))
			isDaemon, _ := cmd.Flags().GetBool("daemon")
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
		execute.Daemon(worker)
	} else {
		log.Info().Msg("Executing in interactive mode ...")
		execute.Interactive(worker)
	}
}

//func worker(ctx context.Context) error {
func worker() {
LOOP:
	for {
		// Calling Sleep method
		time.Sleep(5 * time.Second)
		select {
		case <-execute.Done:
			log.Info().Msg("Graceful termination")
			os.Exit(0)
		case <-execute.Stop:
			log.Warn().Msg("Process terminated by external signal")
			break LOOP		
		case <-execute.Reload:
			log.Info().Msg("Reloading configuration")
		default:
			log.Info().Msg("Execution time ...")
		}
	}
	os.Exit(1)
}
