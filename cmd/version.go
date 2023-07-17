package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   VersionAction,
	Short: "Prints the current version and exits",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n", ITZVersionString)
		if err != nil {
			return fmt.Errorf("error getting version for itz: %v", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
