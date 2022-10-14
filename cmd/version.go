/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the current version and exits",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", ATKVersionString)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
