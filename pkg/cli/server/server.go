package server

import (
	"context"
	"fmt"
	//"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	//"strconv"

	"github.com/rs/zerolog/log"
	gd "github.com/sevlyar/go-daemon"
	"github.com/spf13/cobra"
	//"devcircus.com/cerberus/pkg/config"
)

const defaultTick = 60 * time.Second

var (
	stop = make(chan struct{})
	done = make(chan struct{})
)

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
		if err := runAsDaemon(pidFile); err != nil {
			log.Fatal().Err(err).Msg(err.Error())
		}
	} else {
		log.Info().Msg("Executing in interactive mode ...")
		if err := runAsInteractive(); err != nil {
			log.Fatal().Err(err).Msg(err.Error())
		}
	}
}

func runAsDaemon(pidFile int) error {
	gd.SetSigHandler(termHandler, syscall.SIGTERM)

	//dmn := &gd.Context{
	//	//PidFileName: pidFile,"/tmp/daemonize.pid"
	//	PidFileName: "./",
	//	PidFilePerm: 0644,
	//	WorkDir:     "/",
	//	Umask:       027,
	//}

	dmn := new(gd.Context)

	child, err := dmn.Reborn()
	if err != nil {
		log.Error().Msg(fmt.Sprintf("An error occured while trying to reborn daemon %s", err.Error()))
	}
	if child != nil {
		return err
	}

	defer dmn.Release()

	go worker()
	go func() {
		for {
			time.Sleep(time.Second)
			if _, ok := <-stop; ok {
				log.Info().Msg("Terminating daemon")
			}
		}
	}()

	err = gd.ServeSignals()
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
		return err
	}
	return nil
}

func termHandler(sig os.Signal) error {
	stop <- struct{}{}
	if sig == syscall.SIGTERM {
		<-done
	}
	return gd.ErrStop
}

func runAsInteractive() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()

	go func() {
		for {
			select {
			case <-signalChan:
				log.Info().Msg("Got SIGINT/SIGTERM, exiting.")
				cancel()
				os.Exit(1)
			case <-ctx.Done():
				log.Printf("Done.")
				os.Exit(1)
			}
		}
	}()

	//if err := worker(ctx); err != nil {
	if err := worker(); err != nil {
		return err
	}
	return nil
}

//func worker(ctx context.Context) error {
func worker() error {
	for {
		log.Info().Msg("Execution time ...")
		// Calling Sleep method
		time.Sleep(5 * time.Second)
		//select {
		//case <-ctx.Done():
		//	return nil
		//case <-time.Tick(defaultTick):
		//	resp, err := http.Get("http://www.google.es")
		//	if err != nil {
		//		return err
		//	}
		//	log.Info().Msg(fmt.Sprintf("Status code 200, got: %d", resp.StatusCode))
		//}
	}
}
