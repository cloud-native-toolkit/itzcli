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

var reservationID string

// getReservationCmd represents the viewReservation command
var getReservationCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a specific reservation.",
	Long:  `Get the details of a reservation.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debugf("Getting reservation info for reservation %s...", reservationID)
		return getReservation(cmd, args)
	},
}

func getReservation(cmd *cobra.Command, args []string) error {
	url := viper.GetString("reservation.api.url")
	token := viper.GetString("reservations.api.token")

	if len(url) == 0 {
		return fmt.Errorf("no API url specified for reservation")
	}

	if len(token) == 0 {
		return fmt.Errorf(`There is no API token specified for reservation. Please run the "itz auth login" command to login to TechZone.`)
	}

	if reservationID == "000000000" {
		return fmt.Errorf("no reservation id specified, use --reservation-id ######")
	}

	url = url + reservationID

	logger.Debugf("Using API URL \"%s\" and token \"%s\" to get reservation info...",
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
	} else if len(data) == 0 {
		return fmt.Errorf("no reservation data retrieved, confirm --reservation-id is correct")
	}
	jsoner := reservations.NewJsonReader()
	dataR := bytes.NewReader(data)
	rez, err := jsoner.Read(dataR)
	outer := reservations.NewWriter(jsonFormat)
	
	return outer.WriteOne(reservationCmd.OutOrStdout(), rez)
}

func init() {
	reservationCmd.AddCommand(getReservationCmd)
	getReservationCmd.Flags().StringVar(&reservationID, "reservation-id", "000000000", "Specifies the reservation you want to find information about.")
}
