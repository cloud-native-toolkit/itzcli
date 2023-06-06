package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   VersionAction,
	Short: "Prints the current version and exits",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n", ITZVersionString)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
