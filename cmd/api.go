/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Starts API server http://localhost:8795",
	Long: `Starts API server localhost:8795. Exposes certain commands as HTTP API's.`,
}

func init() {
	RootCmd.AddCommand(apiCmd)
}
