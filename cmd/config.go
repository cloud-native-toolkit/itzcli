package cmd

import (
	"github.com/spf13/cobra"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configures the atk command",
	Long:  `Configures the atk command.`,
}

func init() {
	RootCmd.AddCommand(configureCmd)
}
