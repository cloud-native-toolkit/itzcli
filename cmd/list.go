package cmd

import (
	"fmt"
	"github.com/cloud-native-toolkit/itzcli/pkg/configuration"
	"github.com/cloud-native-toolkit/itzcli/pkg/solutions"
	"github.com/cloud-native-toolkit/itzcli/pkg/techzone"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"reflect"
)

var createdOnly bool
var listAll bool
var componentName string
var owner []string

// listCmd represents the version command
var listCmd = &cobra.Command{
	Use:   ListAction,
	Short: "Lists the summaries of the requested objects",
}

var listReservationCmd = &cobra.Command{
	Use:    pluralOf(ReservationResource),
	Short:  "Displays a list of your current reservations.",
	Long:   `Displays a list of your current IBM Technology Zone reservations.`,
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debug("Listing your reservations...")
		apiConfig, err := LoadApiClientConfig(configuration.TechZone)
		if err != nil {
			return err
		}
		svc, err := techzone.NewReservationWebServiceClient(apiConfig)
		if err != nil {
			return errors.Wrap(err, "could not create web service client")
		}
		w := techzone.NewModelWriter(reflect.TypeOf(techzone.Reservation{}).Name(), GetFormat(cmd))
		sol, err := svc.GetAll(techzone.NoFilter())
		if err != nil {
			return err
		}
		w.WriteMany(cmd.OutOrStdout(), sol)
		return nil
	},
}

var listBuildsCmd = &cobra.Command{
	Use:    pluralOf(BuildResource),
	Short:  fmt.Sprintf("Displays a list of the available %s from the %s catalog.", pluralOf(BuildResource), TechZoneShort),
	Long:   `Displays a list of the available IBM Technology Zone builds.`,
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debugf("Listing the %s %s...", TechZoneFull, BuildResource)
		return listComponents(cmd, args)
	},
}

var listPipelinesCmd = &cobra.Command{
	Use:    pluralOf(PipelineResource),
	Short:  fmt.Sprintf("Displays a list of the available %s from the %s catalog.", pluralOf(PipelineResource), TechZoneShort),
	Long:   `Displays a list of the available IBM Technology Zone pipelines.`,
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debugf("Listing the %s %s...", TechZoneFull, PipelineResource)
		return listComponents(cmd, args)
	},
}

var listEnvironmentCmd = &cobra.Command{
	Use:    pluralOf(EnvironmentResource),
	Short:  fmt.Sprintf("Displays a list of the available %s %s.", TechZoneShort, pluralOf(EnvironmentResource)),
	Long:   `Displays a list of the available IBM Technology Zone environments.`,
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debugf("Listing the %s environments...", TechZoneFull)
		// List component entities.

		apiConfig, err := LoadApiClientConfig(configuration.TechZone)
		if err != nil {
			return err
		}
		svc, err := techzone.NewEnvironmentWebServiceClient(apiConfig)
		if err != nil {
			return errors.Wrap(err, "could not create web service client")
		}
		w := techzone.NewModelWriter(reflect.TypeOf(techzone.Environment{}).Name(), GetFormat(cmd))
		sol, err := svc.GetAll(techzone.NoFilter())
		if err != nil {
			return err
		}
		w.WriteMany(cmd.OutOrStdout(), sol)
		return nil
	},
}

func listComponents(cmd *cobra.Command, args []string) error {
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
	listReservationCmd.Flags().BoolVarP(&listAll, "all", "a", false, "If true, list all reservations (including scheduled)")

	listBuildsCmd.Flags().BoolVarP(&createdOnly, "created", "c", false, "If true, limits the builds to my (created) builds")
	listBuildsCmd.Flags().StringVarP(&componentName, "name", "n", "", "The name of the build")
	listBuildsCmd.Flags().StringSliceVarP(&owner, "owner", "o", owner, "The owner of the build")

	listPipelinesCmd.Flags().BoolVarP(&createdOnly, "created", "c", false, "If true, limits the pipelines to my (created) pipelines")
	listPipelinesCmd.Flags().StringVarP(&componentName, "name", "n", "", "The name of the pipeline")
	listPipelinesCmd.Flags().StringSliceVarP(&owner, "owner", "o", owner, "The owner of the pipeline")

	listCmd.AddCommand(listReservationCmd)
	listCmd.AddCommand(listEnvironmentCmd)
	listCmd.AddCommand(listBuildsCmd)
	listCmd.AddCommand(listPipelinesCmd)
	// Add list
	rootCmd.AddCommand(listCmd)
	// Add the objects to the list command
}
