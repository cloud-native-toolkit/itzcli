package cmd

import (
	"github.com/spf13/cobra"
)

// reservationCmd represents the reservation command
var reservationCmd = &cobra.Command{
	Use:   "reservation",
	Short: "List and get IBM Technology Zone reservations.",
	Long:  `List and get IBM Technology Zone reservations.`,
}

func init() {
	RootCmd.AddCommand(reservationCmd)
}
