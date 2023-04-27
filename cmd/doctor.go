package cmd

import (
	"fmt"
	"github.com/cloud-native-toolkit/itzcli/cmd/dr"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var fixDoctorIssues bool = false

// doctorCmd represents the doctor command
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Checks the environment and configuration",
	Long: `If using the init sub-command, the doctor command will initialize the
environment for first run.
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
}

func init() {
	RootCmd.AddCommand(doctorCmd)
	doctorCmd.Flags().BoolVarP(&fixDoctorIssues, "auto-fix", "f", false, "If true, makes best attempt to fix the issues")
}

func RunDoctor(fix bool) error {
	configChecks := dr.AllConfigChecks
	fileChecks := dr.FileChecks
	allChecks := append(dr.ActionChecks, configChecks...)
	allChecks = append(allChecks, fileChecks...)
	errs := dr.DoChecks(allChecks, fix)
	if len(errs) > 0 {
		logger.Error("One or more requirements unmet; consider using doctor --auto-fix or doctor init to try to resolve them")
		return fmt.Errorf("found %d errors", len(errs))
	}
	return nil
}
