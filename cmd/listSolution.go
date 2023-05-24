package cmd

import (
	"github.com/cloud-native-toolkit/itzcli/pkg/configuration"
	"github.com/cloud-native-toolkit/itzcli/pkg/solutions"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var solutionName string
var owner []string

// listSolutionCmd represents the listReservation command
var listSolutionCmd = &cobra.Command{
	Use:    "list",
	PreRun: SetLoggingLevel,
	Short:  "Lists your IBM Technology Zone solutions.",
	Long:   `Lists the solutions for your IBM Technology Zone user.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debug("Listing your solutions...")
		return listSolutions(cmd, args)
	},
}

var createdOnly bool

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func listSolutions(cmd *cobra.Command, args []string) error {

	filters := solutions.NewFilter(
		solutions.OwnerFilter(owner),
		solutions.KindFilter([]string{"Asset", "Component", "Product"}),
	)
	logger.Debugf("Using filter(s) %s", filters)
	// List component entities.

	apiConfig, err := LoadApiClientConfig(configuration.Backstage)
	if err != nil {
		return err
	}
	svc, err := solutions.NewWebServiceClient(apiConfig)
	if err != nil {
		return errors.Wrap(err, "could not create web service client")
	}
	w := solutions.NewSolutionWriter(GetFormat(cmd))
	sol, err := svc.GetAll(filters)
	if err != nil {
		return err
	}
	w.WriteMany(cmd.OutOrStdout(), sol)
	return nil
}

func init() {
	solutionCmd.AddCommand(listSolutionCmd)
	listSolutionCmd.Flags().BoolVarP(&createdOnly, "created", "c", false, "If true, limits the solutions to my (created) solutions.")
	listSolutionCmd.Flags().StringVarP(&solutionName, "name", "n", "", "The name of the solution")
	listSolutionCmd.Flags().StringSliceVarP(&owner, "owner", "o", owner, "The owner of the solution")
}
