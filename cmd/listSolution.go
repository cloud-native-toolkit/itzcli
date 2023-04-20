package cmd

import (
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/cloud-native-toolkit/itzcli/pkg/solutions"
	"github.com/spf13/cobra"
	"context"
	"go.einride.tech/backstage/catalog"
)

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
	// Create a Software Catalog API client.
	client := catalog.NewClient(
		catalog.WithBaseURL(url),
	)
	filter := solutions.NewFilter(
		solutions.OwnerFilter(owner),
		solutions.KindFilter([]string{"Asset", "Collection", "Product"}),
	).BuildFilter()
	logger.Debugf("Using filter(s) %s", filter)
	// List component entities.
	sols, err := client.ListEntities(context.Background(), &catalog.ListEntitiesRequest{
		Filters: filter,
	})
	if err != nil {
		return err
	}    
	outer := solutions.NewTextWriter()
	for _, entity := range sols.Entities {
		// Standard fields are parsed into Go structs.
		outer.Write(solutionCmd.OutOrStdout(), entity)
	}
	return nil
}

func init() {
	solutionCmd.AddCommand(listSolutionCmd)
	listSolutionCmd.Flags().BoolVarP(&createdOnly, "created", "c", false, "If true, limits the solutions to my (created) solutions.")
}
