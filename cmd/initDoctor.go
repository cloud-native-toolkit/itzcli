package cmd

import (
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:    "init",
	Short:  "Alias for auto-fix, but also quiet.",
	PreRun: SetQuietLogging,
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunDoctor(true)
	},
}

func init() {
	doctorCmd.AddCommand(initCmd)
}
