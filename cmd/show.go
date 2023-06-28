package cmd

import (
	"fmt"
	"github.com/cloud-native-toolkit/itzcli/pkg/techzone"
	"reflect"

	"github.com/cloud-native-toolkit/itzcli/pkg/configuration"
	"github.com/cloud-native-toolkit/itzcli/pkg/solutions"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var buildID string
var pipelineID string
var reservationID string

// showCmd represents the version command
var showCmd = &cobra.Command{
	Use:   ShowAction,
	Short: "Shows the details of the requested single object",
	Long:  `Shows the details of the requested single object.`,
}

var showReservationCmd = &cobra.Command{
	Use:    ReservationResource,
	Short:  "Shows the details of the specific reservation",
	Long:   `Shows the details of the specific IBM Technology Zone reservation.`,
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debug("Listing your reservations...")
		if len(reservationID) == 0 {
			return fmt.Errorf("the reservation id is empty; use --reservation-id to specify the ID of a reservation")
		}
		apiConfig, err := LoadApiClientConfig(configuration.TechZone)
		if err != nil {
			return err
		}
		svc, err := techzone.NewReservationWebServiceClient(apiConfig)
		if err != nil {
			return errors.Wrap(err, "could not create web service client")
		}
		w := techzone.NewModelWriter(reflect.TypeOf(techzone.Reservation{}).Name(), GetFormat(cmd))
		rez, err := svc.Get(reservationID)
		if err != nil {
			return err
		}
		return w.WriteOne(cmd.OutOrStdout(), rez)
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
		return nil
	},
}

var showEnvironmentCmd = &cobra.Command{
	Use:    EnvironmentResource,
	Short:  fmt.Sprintf("Shows the details of the %s %s", TechZoneShort, EnvironmentResource),
	Long:   `Shows the details of the IBM Technology Zone environments.`,
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debugf("Listing the %s environments...", TechZoneFull)
		return nil
	},
}

func init() {

	// Add the parameters to the show commands...
	showBuildsCmd.Flags().StringVar(&buildID, "build-id", "", "ID of the build in the catalog")
	showPipelinesCmd.Flags().StringVar(&pipelineID, "pipeline-id", "", "ID of the build in the catalog")
	showReservationCmd.Flags().StringVar(&reservationID, "reservation-id", "", "ID of the reservation")

	showCmd.AddCommand(showReservationCmd)
	showCmd.AddCommand(showEnvironmentCmd)
	showCmd.AddCommand(showBuildsCmd)
	showCmd.AddCommand(showPipelinesCmd)

	rootCmd.AddCommand(showCmd)
}
