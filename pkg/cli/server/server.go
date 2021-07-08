package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	//"strconv"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	//"devcircus.com/cerberus/pkg/config"
)

const defaultTick = 60 * time.Second

// NewCmdServer start the server
func NewCmdServer() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "start the system server",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info().Msg(fmt.Sprintf("Starting server with pid: %d", os.Getpid()))

			//var port = config.System.Port
			//if port == 0 {
			//	//Default port
			//	http.ListenAndServe(":7321", nil)
			//} else {
			//	//if port is mentioned in config file
			//	http.ListenAndServe(":"+strconv.Itoa(port), nil)
			//}

			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)

			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)

			defer func() {
				signal.Stop(signalChan)
				cancel()
			}()

			go func() {
				for {
					select {
					case s := <-signalChan:
						switch s {
						case os.Interrupt:
							log.Info().Msg("Got SIGINT/SIGTERM, exiting.")
							cancel()
							os.Exit(1)
						}
					case <-ctx.Done():
						log.Printf("Done.")
						os.Exit(1)
					}
				}
			}()

			if err := run(ctx); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
		},
	}
}

func run(ctx context.Context) error {

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.Tick(defaultTick):
			resp, err := http.Get("http://www.google.es")
			if err != nil {
				return err
			}

			log.Info().Msg(fmt.Sprintf("Status code 200, got: %d", resp.StatusCode))

			//if s := resp.Header.Get("server"); s != c.server {
			//	log.Printf("Server header mismatch, got: %s\n", s)
			//}

			//if ct := resp.Header.Get("content-type"); ct != c.contentType {
			//	log.Printf("Content-Type header mismatch, got: %s\n", ct)
			//}

			//if ua := resp.Header.Get("user-agent"); ua != c.userAgent {
			//	log.Printf("User-Agent header mismatch, got: %s\n", ua)
			//}
		}
	}
}
