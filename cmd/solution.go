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
working with the IBM Technology Zone Accelerator Toolkit solutions.

See https://builder.cloudnativetoolkit.dev/ for more information.`,
}

func init() {
	RootCmd.AddCommand(solutionCmd)
}
