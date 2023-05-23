package cmd

import (
	"context"
	"fmt"

	"github.com/cloud-native-toolkit/itzcli/pkg/solutions"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tdabasinskas/go-backstage/v2/backstage"
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
	url := viper.GetString("backstage.api.url")
	if len(url) == 0 {
		return fmt.Errorf("no url specified for backstage")
	}
	logger.Debugf("Using url %s", url)
	// Create a Software Catalog API client.
	c, _ := backstage.NewClient(url, "default", nil)
	filters := solutions.NewFilter(
		solutions.OwnerFilter(owner),
		solutions.KindFilter([]string{"Asset", "Component", "Product"}),
	).BuildFilter()
	logger.Debugf("Using filter(s) %s", filters)
	// List component entities.
	sols, _, err := c.Catalog.Entities.List(context.Background(), &backstage.ListEntityOptions{
		Filters: filters,
	})
	if err != nil {
		return err
	}
	outer := solutions.NewWriter(jsonFormat)
	// Standard fields are parsed into Go structs.
	outer.Write(solutionCmd.OutOrStdout(), sols)
	return nil
}

func init() {
	solutionCmd.AddCommand(listSolutionCmd)
	listSolutionCmd.Flags().BoolVarP(&createdOnly, "created", "c", false, "If true, limits the solutions to my (created) solutions.")
	listSolutionCmd.Flags().StringVarP(&solutionName, "name", "n", "", "The name of the solution")
	listSolutionCmd.Flags().StringSliceVarP(&owner, "owner", "o", owner, "The owner of the solution")
}
