package cmd

import (
	logger "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// listSolutionCmd represents the listReservation command
var listSolutionCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists your TechZone solutions.",
	Long:  `Lists the solutions for your TechZone user.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Listing your solutions...")
	},
}

func init() {
	solutionCmd.AddCommand(listSolutionCmd)
}
