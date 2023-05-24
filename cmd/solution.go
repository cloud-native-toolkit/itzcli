/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// solutionCmd represents the project command
var solutionCmd = &cobra.Command{
	Use:    "solution",
	PreRun: SetLoggingLevel,
	Short:  "Lists metadata, builds, and deploys solutions",
	Long: `The solution command provides a CLI for maintaining
working with components in the IBM Technology Zone catalog.

See https://catalog.techzone.ibm.com for more information.`,
}

func init() {
	RootCmd.AddCommand(solutionCmd)
}
