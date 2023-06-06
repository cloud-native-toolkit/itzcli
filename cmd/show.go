package cmd

import (
	"fmt"

	"github.com/cloud-native-toolkit/itzcli/pkg/configuration"
	"github.com/cloud-native-toolkit/itzcli/pkg/solutions"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var buildID string
var pipelineID string

// showCmd represents the version command
var showCmd = &cobra.Command{
	Use:   ShowAction,
	Short: "Shows the details of the requested single object",
	Long:  `Shows the details of the requested single object.`,
}

var showReservationCmd = &cobra.Command{
	Use:    pluralOf(ReservationResource),
	Short:  "Shows the details of the specific reservation",
	Long:   `Shows the details of the specific IBM Technology Zone reservation.`,
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debug("Listing your reservations...")
		return listReservations(cmd, args)
	},
}

var showBuildsCmd = &cobra.Command{
	Use:    BuildResource,
	Short:  fmt.Sprintf("Shows the details of the specific %s from the %s catalog", BuildResource, TechZoneShort),
	Long:   `Shows the details of the IBM Technology Zone builds.`,
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Info("Getting your solution...")
		if len(buildID) == 0 {
			return fmt.Errorf("solution id is empty")
		}
		apiConfig, err := LoadApiClientConfig(configuration.Backstage)
		if err != nil {
			return err
		}
		svc, err := solutions.NewWebServiceClient(apiConfig)
		if err != nil {
			return errors.Wrap(err, "could not create web service client")
		}
		w := solutions.NewSolutionWriter(GetFormat(cmd))
		sol, err := svc.Get(buildID)
		if err != nil {
			return err
		}
		w.Write(cmd.OutOrStdout(), sol)
		return nil
	},
}

var showPipelinesCmd = &cobra.Command{
	Use:    PipelineResource,
	Short:  fmt.Sprintf("Shows the details of the specific %s from the %s catalog", PipelineResource, TechZoneShort),
	Long:   `Shows the details of the IBM Technology Zone pipelines.`,
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debugf("Listing the %s %s...", TechZoneFull, PipelineResource)
		return listReservations(cmd, args)
	},
}

var showEnvironmentCmd = &cobra.Command{
	Use:    EnvironmentResource,
	Short:  fmt.Sprintf("Shows the details of the %s %s", TechZoneShort, EnvironmentResource),
	Long:   `Shows the details of the IBM Technology Zone environments.`,
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debugf("Listing the %s environments...", TechZoneFull)
		return listReservations(cmd, args)
	},
}

func init() {

	// Add the parameters to the show commands...
	showBuildsCmd.Flags().StringVar(&buildID, "build-id", "", "ID of the build in the catalog")
	showPipelinesCmd.Flags().StringVar(&pipelineID, "pipeline-id", "", "ID of the build in the catalog")

	showCmd.AddCommand(showReservationCmd)
	showCmd.AddCommand(showEnvironmentCmd)
	showCmd.AddCommand(showBuildsCmd)
	showCmd.AddCommand(showPipelinesCmd)

	rootCmd.AddCommand(showCmd)
}
