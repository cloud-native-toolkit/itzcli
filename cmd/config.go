package cmd

import (
	"github.com/spf13/cobra"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configures the itz command",
	Long:  `Configures the itz command.`,
}

func init() {
	RootCmd.AddCommand(configureCmd)
}
