package cmd

import (
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var reserveCmd = &cobra.Command{
	Use:    ReserveAction,
	Short:  "Allows you to reserve environments",
	Long:   "Allows you to reserve environments",
	PreRun: SetLoggingLevel,
	// TODO: This command is temporarily not supported, but this is a placeholder
	// for future versions.
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debug("Reserving your environment...")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(reserveCmd)
}
