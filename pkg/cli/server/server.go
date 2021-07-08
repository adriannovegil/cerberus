package server

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ()

// NewCmdServer start the server
func NewCmdServer() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "start the system server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Run server ...")
		},
	}
}
