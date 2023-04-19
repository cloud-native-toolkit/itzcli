package cmd

import (
	"bytes"
	"fmt"

	"github.com/cloud-native-toolkit/itzcli/pkg"
	"github.com/cloud-native-toolkit/itzcli/pkg/reservations"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listAllRez bool

const ReservationListPermissionsError = `
Permissions error while trying to read from your list of reservations. The most
common cause is an expired or bad API token. You can resolve this issue by going
to https://techzone.ibm.com/my/profile to get your API token, save it in a file
(e.g., /path/to/token.txt) and use the command:

    $ itz auth login --from-file /path/to/token.txt --service reservations

`

// listReservationCmd represents the listReservation command
var listReservationCmd = &cobra.Command{
	Use:    "list",
	Short:  "Lists your current IBM Technology Zone reservations.",
	Long:   `Lists your current IBM Technology Zone reservations.`,
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debug("Listing your reservations...")
		return listReservations(cmd, args)
	},
}

func listReservations(cmd *cobra.Command, args []string) error {
	url := viper.GetString("reservations.api.url")
	token := viper.GetString("reservations.api.token")

	if len(url) == 0 {
		return fmt.Errorf("no API url specified for reservations")
	}

	if len(token) == 0 {
		return fmt.Errorf("no API token specified for reservation. Please run the auth login command to login to TechZone.")
	}

	logger.Debugf("Using API URL \"%s\" and token \"%s\" to get list of reservations...",
		url, token)

	data, err := pkg.ReadHttpGetTWithFunc(url, token, func(code int) error {
		logger.Debugf("Handling HTTP return code %d...", code)
		if code == 401 {
			pkg.WriteMessage(ReservationListPermissionsError, reservationCmd.OutOrStdout())
		}
		return nil
	})
	if err != nil {
		return err
	}
	jsoner := reservations.NewJsonReader()
	dataR := bytes.NewReader(data)
	rez, err := jsoner.ReadAll(dataR)

	logger.Debugf("Found %d reservations.", len(rez))
	outer := reservations.NewTextWriter()
	matches := 0
	if listAllRez {
		// --list-all includes the statuses, plus deleted.
		matches, err  =  outer.WriteFilter(reservationCmd.OutOrStdout(), rez, reservations.FilterByStatusSlice([]string{"Ready", "Scheduled", "Provisioning", "Deleted"}))
	} else {
		matches, err = outer.WriteFilter(reservationCmd.OutOrStdout(), rez, reservations.FilterByStatusSlice([]string{"Ready", "Scheduled", "Provisioning"}))
		// check to see if we had any matches for our active reservation filter
		// if not, then print a error to the user
		if (matches == 0) {
			logger.Print("No reservations found. Use -a to view expired reservations")
		}
	}
	return err 
}

func init() {
	reservationCmd.AddCommand(listReservationCmd)
	listReservationCmd.Flags().BoolVarP(&listAllRez, "all", "a", false, "If true, list all reservations (including scheduled)")
}
