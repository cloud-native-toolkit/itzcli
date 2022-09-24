package cmd

/*
Catalog api:
 Prod: https://api.techzone.ibm.com/swagger/
 Stage: https://techzone-staging-api.dal1a.ciocloud.nonprod.intranet.ibm.com/swagger
 Test: https://techzone-test-api.dal1a.ciocloud.nonprod.intranet.ibm.com/swagger/

Reservations:
 Prod: https://reservations.techzone.ibm.com/swagger/
 Stage: https://techzone-staging-reservations.dal1a.ciocloud.nonprod.intranet.ibm.com/swagger
 Test:     http://techzone-test-reservations.dal1a.ciocloud.nonprod.intranet.ibm.com/swagger

Auth:
 Prod: https://auth.techzone.ibm.com/swagger/
 Stage: https://techzone-staging-auth.dal1a.ciocloud.nonprod.intranet.ibm.com/swagger

Journal/Metrics:
 Prod: https://accounting.techzone.ibm.com/swagger/
 Stage: https://techzone-staging-accounting.dal1a.ciocloud.nonprod.intranet.ibm.com/swagger
*/

import (
	"bytes"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.ibm.com/skol/atkcli/pkg"
	"github.ibm.com/skol/atkcli/pkg/reservations"

	"github.com/spf13/cobra"
)

// listReservationCmd represents the listReservation command
var listReservationCmd = &cobra.Command{
	Use:    "list",
	Short:  "Lists your current TechZone reservations.",
	Long:   `Lists your current TechZone reservations.`,
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Info("Listing your reservations...")
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
		return fmt.Errorf("no API token specified for reservations")
	}

	logger.Debugf("Using API URL \"%s\" and token \"%s\" to get list of reservations...",
		url, token)

	data, err := pkg.ReadHttpGetT(url, token)
	if err != nil {
		return err
	}
	jsoner := reservations.NewJsonReader()
	dataR := bytes.NewReader(data)
	rez, err := jsoner.ReadAll(dataR)

	logger.Debugf("Found %d reservations.", len(rez))
	outer := reservations.NewTextWriter()
	outer.WriteFilter(reservationCmd.OutOrStdout(), rez, reservations.FilterByStatus("Ready"))

	return nil
}

func init() {
	reservationCmd.AddCommand(listReservationCmd)
	// Here you will define your flags and configuration settings.
}
