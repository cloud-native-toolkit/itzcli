package cmd

import (
	"bytes"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.ibm.com/skol/atkcli/pkg"
	"github.ibm.com/skol/atkcli/pkg/solutions"
	"strings"

	"github.com/spf13/cobra"
)

// listSolutionCmd represents the listReservation command
var listSolutionCmd = &cobra.Command{
	Use:    "list",
	PreRun: SetLoggingLevel,
	Short:  "Lists your TechZone solutions.",
	Long:   `Lists the solutions for your TechZone user.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debug("Listing your solutions...")
		return listSolutions(cmd, args)
	},
}

func listSolutions(cmd *cobra.Command, args []string) error {
	// HACK: This will eventually be a URL and not a URL or a file path.
	// Load up the reader based on the URI provided for the solution
	uri := viper.GetString("builder.api.url")
	token := viper.GetString("builder.api.token")

	if len(uri) == 0 {
		return fmt.Errorf("no API url specified for builder")
	}

	if len(token) == 0 {
		return fmt.Errorf("no API token specified for builder")
	}

	var data []byte
	var err error
	if strings.HasPrefix(uri, "https://") || strings.HasPrefix(uri, "http://") {
		logger.Debugf("Using API URL \"%s\" and token \"%s\" to get list of reservations...", uri, token)
		data, err = pkg.ReadHttpGetT(uri, token)
	} else {
		logger.Debugf("Loading solutions from file: \"%s\"", uri)
		data, err = pkg.ReadFile(uri)
	}

	if err != nil {
		return err
	}
	jsoner := solutions.NewJsonReader()
	dataR := bytes.NewReader(data)
	sols, err := jsoner.ReadAll(dataR)

	logger.Debugf("Found %d reservations.", len(sols))
	outer := solutions.NewTextWriter()
	outer.WriteAll(solutionCmd.OutOrStdout(), sols)

	return nil

	return nil
}

func init() {
	solutionCmd.AddCommand(listSolutionCmd)
}
