package cmd

import (
	logger "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// getReservationCmd represents the viewReservation command
var getReservationCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a specific reservation.",
	Long:  `Get the details of a reservation.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("Getting your reservation...")
	},
}

func init() {
	reservationCmd.AddCommand(getReservationCmd)
}
