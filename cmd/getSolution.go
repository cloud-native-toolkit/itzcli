package cmd

// import (
// 	"fmt"

// 	"github.com/cloud-native-toolkit/itzcli/pkg/configuration"
// 	"github.com/cloud-native-toolkit/itzcli/pkg/solutions"
// 	"github.com/pkg/errors"
// 	logger "github.com/sirupsen/logrus"
// 	"github.com/spf13/cobra"
// )

// var solutionID string

// // getSolutionCmd is the command for getting information about a single solution
// var getSolutionCmd = &cobra.Command{
// 	Use:   "get",
// 	Short: "Gets details for a specific solution.",
// 	Long:  `Get the details of a solution.`,
// 	RunE: func(cmd *cobra.Command, args []string) error {
// 		logger.Info("Getting your solution...")
// 		if len(solutionID) == 0 {
// 			return fmt.Errorf("solution id is empty")
// 		}
// 		apiConfig, err := LoadApiClientConfig(configuration.Backstage)
// 		if err != nil {
// 			return err
// 		}
// 		svc, err := solutions.NewWebServiceClient(apiConfig)
// 		if err != nil {
// 			return errors.Wrap(err, "could not create web service client")
// 		}
// 		w := solutions.NewSolutionWriter(GetFormat(cmd))
// 		sol, err := svc.Get(solutionID)
// 		if err != nil {
// 			return err
// 		}
// 		w.Write(cmd.OutOrStdout(), sol)
// 		return nil
// 	},
// }

// func init() {
// 	solutionCmd.AddCommand(getSolutionCmd)
// 	getSolutionCmd.Flags().StringVar(&solutionID, "solution",
// 		"",
// 		"Specifies the name or UID of the solution to get from the catalog")
// }
