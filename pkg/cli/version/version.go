package version

import (
	"fmt"

	"github.com/spf13/cobra"

	"devcircus.com/cerberus/pkg/config/system"
)

var ()

// NewCmdVersion prints version
func NewCmdVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "show version statement",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("datanerd", system.GetVersion())
		},
	}
}
