/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// getReservationCmd represents the viewReservation command
var getReservationCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a specific reservation.",
	Long:  `Get the details of a reservation.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("viewReservation called")
	},
}

func init() {
	reservationCmd.AddCommand(getReservationCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getReservationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getReservationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
