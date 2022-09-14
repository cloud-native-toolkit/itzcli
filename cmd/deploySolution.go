/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	logger "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// deploySolutionCmd represents the deployProject command
var deploySolutionCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys the specified solution.",
	Long: `Use this command to deploy the specified solution
locally in your own environment.
`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Publishing your solution...")
	},
}

func init() {
	solutionCmd.AddCommand(deploySolutionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deploySolutionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deploySolutionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
