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
	Use:   "login",
	Short: "Uses your browser to authenticate with TechZone.",
	Long: `
Opens a browser window for you to authenticate with IBM Technology Zone using
your IBMid. 

Upon successful login, the CLI updates the configuration with an authentication
token that will be used to access the IBM Technology Zone API as well as the 
IBM Technology Zone Catalog API.
`,
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
		apiArgs := []string{"execute", "api"}
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
	viper.Set(fmt.Sprintf("%s.api.token", "techzone"), string(token))
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
