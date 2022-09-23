/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.ibm.com/skol/atkcli/pkg"
	"net/url"
)

var fn string
var sol string

const bifrost = "bifrost"
const builder = "ci"

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

	for _, svc := range []string{builder, bifrost} {
		cfg := newApiCfg(svc)
		svcUrl, err := url.Parse(cfg.uri)
		if err != nil {
			return err
		}

		if cfg.isLocal {
			logger.Debugf("Using local agent at <%s> for deployment..", svcUrl)
			err := pkg.StartSvcImg(cfg.img, GetPort(svcUrl))
			if err != nil {
				return err
			}
		} else {
			logger.Debugf("Using service at <%s> for deployment", svcUrl)
		}
	}
	return nil
}

type apiCfg struct {
	isLocal bool
	img     string
	uri     string
}

func newApiCfg(key string) *apiCfg {
	return &apiCfg{
		isLocal: viper.GetBool(fmt.Sprintf("%s.api.local", key)),
		img:     viper.GetString(fmt.Sprintf("%s.api.image", key)),
		uri:     viper.GetString(fmt.Sprintf("%s.api.url", key)),
	}
}

func GetPort(uri *url.URL) string {
	return uri.Port()
}
