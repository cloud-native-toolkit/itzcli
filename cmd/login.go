package cmd

import (
	"fmt"
	"os"
	"os/signal"

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
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			<-c
			logger.Debugln("User exited during auth login...")
			os.Exit(1)
		}()
		// Handle the legacy way of login via the text file
		if filePath != "" {
			return TextFileLogin(cmd, args)
		}
		// start the api
		apiArgs := []string{"api", "start"}
		rootCmd.SetArgs(apiArgs) // set the command's args
		// run the command in the background
		go rootCmd.Execute()
		return auth.GetToken()
	},
}

func TextFileLogin(cmd *cobra.Command, args []string) error {

	logger.Debugf("Saving login credentials for reservations using token in file %s...", filePath)
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
	loginCmd.Flags().StringVarP(&filePath, "from-file", "f", "", "The name of the file that contains the token.")
	rootCmd.AddCommand(loginCmd)
}
