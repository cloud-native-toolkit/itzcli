/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.ibm.com/skol/atkcli/pkg"
	"github.ibm.com/skol/atkmod"
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
	bifrostURL, err := url.Parse(viper.GetString("bifrost.api.url"))
	if err != nil {
		return fmt.Errorf("error trying to parse \"bifrost.api.url\", looks like a bad URL (value was: %s): %v", err, viper.GetString("bifrost.api.url"))
	}
	builderURL, err := url.Parse(viper.GetString("ci.api.url"))
	if err != nil {
		return fmt.Errorf("error trying to parse \"ci.api.url\", looks like a bad URL (value was: %s): %v", err, viper.GetString("ci.api.url"))
	}

	services := []pkg.Service{
		{
			DisplayName: "builder",
			ImgName:     viper.GetString("ci.api.image"),
			IsLocal:     viper.GetBool("ci.api.local"),
			URL:         builderURL,
			PreStart:    pkg.StatusHandler,
			Start:       pkg.StartHandler,
		},
		{
			DisplayName: "integration",
			ImgName:     viper.GetString("bifrost.api.image"),
			IsLocal:     viper.GetBool("bifrost.api.local"),
			URL:         bifrostURL,
			PreStart:    pkg.StatusHandler,
			Start:       pkg.StartHandler,
		},
	}

	out := new(bytes.Buffer)
	ctx := &atkmod.RunContext{
		Out: out,
		Log: *logger.StandardLogger(),
	}

	err = pkg.StartupServices(ctx, services, pkg.Sequential)

	if err != nil {
		return err
	}

	// TODO: Now the services are started, we can use them like we would...

	return nil
}

