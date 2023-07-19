package cmd

import (
	"fmt"
	"reflect"

	"github.com/cloud-native-toolkit/itzcli/pkg/configuration"
	"github.com/cloud-native-toolkit/itzcli/pkg/solutions"
	"github.com/cloud-native-toolkit/itzcli/pkg/techzone"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var createdOnly bool
var listAll bool
var componentNames []string
var owner []string

// listCmd represents the version command
var listCmd = &cobra.Command{
	Use:   ListAction,
	Short: "Lists the summaries of the requested objects",
}

var listReservationCmd = &cobra.Command{
	Use:   pluralOf(ReservationResource),
	Short: "Displays a list of your current reservations.",
	Long: `
Displays a list of your current IBM Technology Zone reservations.

By default, the CLI limits the reservations listed to those in "Pending",
"Provisioning", or "Ready" status. To view reservations in "Deleted" or
"Expired" status, use --all ("-a") to list all of your reservations.

The default output format is text in a tabular list. For scripting or
programmatic interaction, specify the --json flag to the command to view the
output in JSON format.

Examples:

    itz list reservations --json
    itz list reservations --all
`,
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
		var filter techzone.Filter
		if listAll {
			filter = techzone.FilterByStatusSlice([]string{"Deleted", "Expired", "Pending", "Ready", "Provisioning"})
		} else {
			filter = techzone.FilterByStatusSlice([]string{"Pending", "Ready", "Provisioning"})
		}
		sol, err := svc.GetAll(filter)
		if err != nil {
			return err
		}
		w.WriteMany(cmd.OutOrStdout(), sol)
		return nil
	},
}

var listPipelinesCmd = &cobra.Command{
	Use:   pluralOf(PipelineResource),
	Short: fmt.Sprintf("Displays a list of the available %s from the %s catalog.", pluralOf(PipelineResource), TechZoneShort),
	Long: `
Displays a list of the available IBM Technology Zone (TechZone) pipelines from
the catalog.

From the TechZone catalog (see https://catalog.techzone.ibm.com/), a pipline is
a deployable component. It must be of kind "Component" and type "pipeline" to be
deployed to a cluster.

Example:

    itz list pipelines
	itz list pipleines -o user:example.owner@ibm.com
	itz list pipleines -n "Deployer CP4S 1.10"
    itz list pipelines --json
`,
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debugf("Listing the %s %s...", TechZoneFull, PipelineResource)
		filters := solutions.NewFilter(
			solutions.OwnerFilter(owner),
			solutions.ComponentNameFilter(componentNames),
			solutions.KindFilter([]string{"Component"}),
			solutions.TypeFilter([]string{"pipeline"}),
		)
		return listComponents(cmd, filters)
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
	Hidden: true,
}

func listComponents(cmd *cobra.Command, filters *solutions.Filter) error {

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
	listReservationCmd.Flags().BoolVarP(&listAll, "all", "a", false, "If true, list all reservations (including expired)")

	listPipelinesCmd.Flags().StringSliceVarP(&componentNames, "name", "n", componentNames, "The name of the pipeline")
	listPipelinesCmd.Flags().StringSliceVarP(&owner, "owner", "o", owner, "The owner of the pipeline")

	listCmd.AddCommand(listReservationCmd)
	listCmd.AddCommand(listEnvironmentCmd)
	listCmd.AddCommand(listPipelinesCmd)
	// Add list
	rootCmd.AddCommand(listCmd)
	// Add the objects to the list command
}
