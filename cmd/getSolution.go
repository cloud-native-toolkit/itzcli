package cmd

import (
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// getReservationCmd represents the viewReservation command
var getSolutionCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a specific reservation.",
	Long:  `Get the details of a reservation.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Getting your solution...")
	},
}

func init() {
	solutionCmd.AddCommand(getSolutionCmd)
}
