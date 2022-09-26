package cmd

import (
	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage tokens and authentication to APIs.",
	Long:  `Manage tokens and authentication to APIs.`,
}

func init() {
	rootCmd.AddCommand(authCmd)
}
