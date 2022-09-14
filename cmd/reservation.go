package cmd

import (
	"github.com/spf13/cobra"
)

// reservationCmd represents the reservation command
var reservationCmd = &cobra.Command{
	Use:   "reservation",
	Short: "List and get TechZone reservations.",
	Long:  `List and get TechZone reservations.`,
}

func init() {
	rootCmd.AddCommand(reservationCmd)
}
