/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/cloud-native-toolkit/itzcli/api"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts API server http://localhost:8795",
	Long: `Starts API server localhost:8795. Exposes certain commands as HTTP API's.`,
	Run: func(cmd *cobra.Command, args []string) {
		api.StartServer(RootCmd)
	},
}

func init() {
	apiCmd.AddCommand(startCmd)
}
