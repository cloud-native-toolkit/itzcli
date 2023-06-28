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
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debug("Reserving your environment...")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(reserveCmd)
}
