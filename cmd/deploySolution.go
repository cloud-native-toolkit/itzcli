/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.ibm.com/Nathan-Good/atkcli/pkg"
	"net/url"
)

var fn string
var sol string

// deploySolutionCmd represents the deployProject command
var deploySolutionCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys the specified solution.",
	Long: `Use this command to deploy the specified solution
locally in your own environment.
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		SetLoggingLevel(cmd, args)
		if len(fn) == 0 && len(sol) == 0 {
			return fmt.Errorf("either \"--solution\" or \"--file\" must be specified.")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Infof("Deploying solution \"%s\"...", sol)
		return DeploySolution(cmd, args)
	},
}

func init() {
	solutionCmd.AddCommand(deploySolutionCmd)
	deploySolutionCmd.Flags().StringVarP(&fn, "file", "f", "", "The full path to the solution file to be deployed.")
	deploySolutionCmd.Flags().StringVarP(&sol, "solution", "s", "", "The name of the solution to be deployed.")
	deploySolutionCmd.MarkFlagsMutuallyExclusive("file", "solution")
}

// DeploySolution deploys the solution by handing it off to the bifrost
// API
func DeploySolution(cmd *cobra.Command, args []string) error {
	// Load up the reader based on the URI provided for the solution
	bifrostUrl := viper.GetString("bifrost.api.url")
	bifrostApi, err := url.Parse(bifrostUrl)
	if err != nil {
		return nil
	}

	if IsLocal() {
		logger.Debugf("Using local agent at <%s> for deployment..", bifrostApi)
		err := pkg.StartUpBifrost(GetPort(bifrostApi))
		if err != nil {
			return err
		}
	} else {
		logger.Debugf("Using service at <%s> for deployment", bifrostApi)
	}

	return nil
}

// IsLocal returns true if the given host string (use a url.Parse()) to get
// it) is the localhost. It does not mind if you give it the port.
func IsLocal() bool {
	return viper.GetBool("bifrost.api.local")
}

func GetPort(uri *url.URL) string {
	return uri.Port()
}
