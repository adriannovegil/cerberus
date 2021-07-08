package server

import (
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"devcircus.com/cerberus/pkg/config"
)

var ()

// NewCmdServer start the server
func NewCmdServer() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "start the system server",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info().Msg("Starting server ...")
			var port = config.System.Port
			if port == 0 {
				//Default port
				http.ListenAndServe(":7321", nil)
			} else {
				//if port is mentioned in config file
				http.ListenAndServe(":"+strconv.Itoa(port), nil)
			}
		},
	}
}
