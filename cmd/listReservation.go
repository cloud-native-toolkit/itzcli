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
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.ibm.com/Nathan-Good/atkcli/pkg"

	"github.com/spf13/cobra"
)

// listReservationCmd represents the listReservation command
var listReservationCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists your current TechZone reservations.",
	Long:  `Lists your current TechZone reservations.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if verbose {
			logger.SetLevel(logger.DebugLevel)
		} else {
			logger.SetLevel(logger.InfoLevel)
		}
		logger.Info("Listing your reservations...")
		return listReservations(cmd, args)
	},
}

func listReservations(cmd *cobra.Command, args []string) error {
	url := viper.GetString("reservations.api.url")
	token := viper.GetString("reservations.api.token")
	logger.Debugf("Using API URL \"%s\" and token \"%s\" to get list of reservations...",
		url, token)

	data, err := pkg.ReadHttpGet(url, token)
	if err != nil {
		return err
	}
	jsoner := pkg.NewJsonReader()
	dataR := bytes.NewReader(data)
	rez, err := jsoner.ReadAll(dataR)

	logger.Debugf("Found %d reservations.", len(rez))
	outer := pkg.NewTextWriter()
	outer.WriteFilter(reservationCmd.OutOrStdout(), rez, pkg.FilterByStatus("Ready"))

	return nil
}

func init() {
	reservationCmd.AddCommand(listReservationCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listReservationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listReservationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
