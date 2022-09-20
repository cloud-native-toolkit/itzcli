package cmd

import (
	logger "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// listSolutionCmd represents the listReservation command
var listSolutionCmd = &cobra.Command{
	Use:    "list",
	PreRun: SetLoggingLevel,
	Short:  "Lists your TechZone solutions.",
	Long:   `Lists the solutions for your TechZone user.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Info("Listing your solutions...")
		return listSolutions(cmd, args)
	},
}

func listSolutions(cmd *cobra.Command, args []string) error {

	return nil
}

func init() {
	solutionCmd.AddCommand(listSolutionCmd)
}
