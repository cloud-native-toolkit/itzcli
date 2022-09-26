package cmd

import (
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.ibm.com/skol/atkcli/pkg"
)

var svcName string
var filePath string

const (
	solutionSvc    string = "builder"
	integrationSvc string = "bifrost"
	rezSvc         string = "reservations"
	builderSvc     string = "ci"
)

var allSvcs = []string{
	solutionSvc,
	integrationSvc,
	rezSvc,
	builderSvc,
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:    "login",
	Short:  "Stores tokens in the configuration for the given service.",
	Long:   `Stores tokens in the configuration for the given service.`,
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !contains(allSvcs, svcName) {
			return fmt.Errorf("service %s not found in supported list of services: %v", svcName, allSvcs)
		}

		logger.Debugf("Saving login credentials for %s using token in file %s...", svcName, filePath)
		token, err := pkg.ReadFile(filePath)
		if err != nil {
			return err
		}
		viper.Set(fmt.Sprintf("%s.api.token", svcName), string(token))
		err = viper.WriteConfig()
		if err != nil {
			return err
		}
		logger.Tracef("Finished writing credentials for %s using token in file %s...", svcName, filePath)
		return nil
	},
}

func init() {
	authCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&svcName, "service-name", "s", "", "The name of the service to login to.")
	loginCmd.Flags().StringVarP(&filePath, "from-file", "f", "", "The name of the file that contains the token.")
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
