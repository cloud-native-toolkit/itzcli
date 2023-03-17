/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var solutionName string
var owner []string
// solutionCmd represents the project command
var solutionCmd = &cobra.Command{
	Use:    "solution",
	PreRun: SetLoggingLevel,
	Short:  "Lists metadata, builds, and deploys solutions",
	Long: `The solution command provides a CLI for maintaining
working with the IBM Technology Zone Accelerator Toolkit solutions.

See https://builder.cloudnativetoolkit.dev/ for more information.`,
}

func init() {
	RootCmd.AddCommand(solutionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	solutionCmd.PersistentFlags().StringVarP(&solutionName, "name", "n", "", "The name of the solution")
	solutionCmd.PersistentFlags().StringArrayVarP(&owner, "owner", "o", owner, "The owner of the solution")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// solutionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
