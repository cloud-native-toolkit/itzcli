package cmd

import (
	"fmt"
	"github.com/cloud-native-toolkit/itzcli/pkg"
	"github.com/cloud-native-toolkit/itzcli/pkg/auth"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)


var filePath string


// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:    "login",
	Short:  "Stores tokens in the configuration for the given service.",
	Long:   `Stores tokens in the configuration for the given service.`,
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Handle the legacy way of login via the text file
		if filePath != "" {
			return TextFileLogin(cmd, args)
		}
		return auth.GetToken()
	},
}

func TextFileLogin(cmd *cobra.Command, args []string) error {
	
	logger.Debugf("Saving login credentials for %s using token in file %s...", "reservations", filePath)
	token, err := pkg.ReadFile(filePath)
	if err != nil {
		return err
	}
	viper.Set(fmt.Sprintf("%s.api.token", "reservations"), string(token))
	err = viper.WriteConfig()
	if err != nil {
		return err
	}
	logger.Tracef("Finished writing credentials for %s using token in file %s...", "reservations", filePath)
	return nil
}

func init() {
	authCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&filePath, "from-file", "f", "", "The name of the file that contains the token.")
}
