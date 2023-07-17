package cmd

import (
	"fmt"

	"github.com/cloud-native-toolkit/itzcli/cmd/dr"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var fixDoctorIssues = false

// doctorCmd represents the doctor command
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Checks the environment and configuration",
	Long: `
Checks the environment and verifies configurations, required binaries, and
creates required directories on the filesystem in the home ("~/.itz") directory.

When using the "--auto-fix" ("-f") flag, the doctor command will initialize the
environment for first run. This means the command will:

* Create the ~/.itz home directory and any required sub-directories.
* Create the "cli-config.yaml" file in the home directory with initial values.
* If a podman machine is running, resolve a bug with the podman machine system
time becoming out of sync.

Examples:

    itz doctor # runs but prints warnings without fixing them
    itz doctor -f # attempts to fix any missing configurations, files, etc.

`,
	PreRun: SetLoggingLevel,
	// Perform the checks on the system to make sure that ITZ is OK to run
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debug("Checking the environment...")
		return RunDoctor(fixDoctorIssues)
	},
	// The usage usually prints if there is an error, but in this case we do not
	// want to print the usage.
	SilenceUsage: true,
	// Hidden:       true,
}

func init() {
	rootCmd.AddCommand(doctorCmd)
	doctorCmd.Flags().BoolVarP(&fixDoctorIssues, "auto-fix", "f", false, "If true, makes best attempt to fix the issues")
}

func RunDoctor(fix bool) error {
	configChecks := dr.AllConfigChecks
	fileChecks := dr.FileChecks
	allChecks := append(fileChecks, dr.ActionChecks...)
	allChecks = append(allChecks, configChecks...)
	errs := dr.DoChecks(allChecks, fix)
	if len(errs) > 0 {
		logger.Error("One or more requirements unmet; consider using doctor --auto-fix to resolve them")
		return fmt.Errorf("found %d errors", len(errs))
	}
	return nil
}
