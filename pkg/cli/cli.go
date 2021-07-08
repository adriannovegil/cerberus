package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"devcircus.com/cerberus/pkg/cli/server"
	"devcircus.com/cerberus/pkg/cli/version"
	"devcircus.com/cerberus/pkg/config"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   config.System.AppName,
	Short: "Cerberus - Monitor your URIs from your any place and get notified when your system is down or response time is more than expected.",
	Long: `Monitor your URIs from your any place and get notified when your system is
down or response time is more than expected.

For more information see the official tool documentation.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// init method
func init() {
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// Adding the commands
	rootCmd.AddCommand(version.NewCmdVersion())
	rootCmd.AddCommand(server.NewCmdServer())
}
