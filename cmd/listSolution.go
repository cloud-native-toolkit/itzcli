package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/cloud-native-toolkit/itzcli/cmd/dr"
	"github.com/cloud-native-toolkit/itzcli/pkg"
	"github.com/cloud-native-toolkit/itzcli/pkg/solutions"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
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
	// HACK: This will eventually be a URL and not a URL or a file path.
	// Load up the reader based on the URI provided for the solution
	uri := viper.GetString("builder.api.url")
	refreshToken := viper.GetString("builder.api.refresh_token")
	refreshToken = strings.TrimSpace(refreshToken)

	if len(uri) == 0 {
		return fmt.Errorf("no API url specified for builder")
	}

	if len(refreshToken) == 0 {
		return fmt.Errorf("could not get refresh token for builder API")
	}

	tokenClient := createTokenRefreshPostClient(uri, refreshToken)
	var tokenResponse tokenResponse
	tokenClient.ResponseHandler = func(reader io.ReadCloser) error {
		defer reader.Close()
		err := json.NewDecoder(reader).Decode(&tokenResponse)
		if err != nil {
			return err
		}
		return nil
	}
	err := pkg.Exec(tokenClient)
	if err != nil {
		return err
	}

	logger.Debugf("Token response: %+v", tokenResponse)
	if len(tokenResponse.AccessToken) > 0 && tokenResponse.ExpiresIn > 0 {
		logger.Debugf("Got token: %s", tokenResponse.AccessToken)
		if refreshToken != tokenResponse.RefreshToken {
			viper.Set("builder.api.refresh_token", tokenResponse.RefreshToken)
			viper.WriteConfig()
		}
	}

	var data []byte
	if strings.HasPrefix(uri, "https://") || strings.HasPrefix(uri, "http://") {
		// Use the refresh_token from the configuration file to get the access_token and the
		// id_token to call the API...
		logger.Debugf("Using API URL \"%s\" and token \"%s\" to get list of solutions...", uri, tokenResponse.AccessToken)
		token := fmt.Sprintf("%s %s", tokenResponse.AccessToken, tokenResponse.IdToken)
		if !createdOnly {
			data, err = pkg.ReadHttpGetTWithFunc(fmt.Sprintf("%s/solutions", uri), token, func(code int) error {
				logger.Debugf("Handling HTTP return code %d...", code)
				if code == 401 {
					pkg.WriteMessage(dr.SolutionsListPermissionsError, reservationCmd.OutOrStdout())
				}
				return nil
			})
		} else {
			username := viper.GetString("builder.api.username")
			data, err = pkg.ReadHttpGetTWithFunc(fmt.Sprintf("%s/users/%s/solutions", uri, username), token, func(code int) error {
				logger.Debugf("Handling HTTP return code %d...", code)
				if code == 401 {
					pkg.WriteMessage(dr.SolutionsListPermissionsError, reservationCmd.OutOrStdout())
				}
				return nil
			})
		}
	} else {
		logger.Debugf("Loading solutions from file: \"%s\"", uri)
		data, err = pkg.ReadFile(uri)
	}

	if err != nil {
		return err
	}
	jsoner := solutions.NewJsonReader()
	dataR := bytes.NewReader(data)
	sols, err := jsoner.ReadAll(dataR)

	logger.Debugf("Found %d reservations.", len(sols))
	outer := solutions.NewTextWriter()
	return outer.WriteAll(solutionCmd.OutOrStdout(), sols)
}

func createTokenRefreshPostClient(bUrl string, t string) *pkg.ServiceClient {
	return &pkg.ServiceClient{
		BaseURL: fmt.Sprintf("%s/token", bUrl),
		Method:  http.MethodPost,
		FParams: func() map[string]string {
			m := make(map[string]string)
			m["refresh_token"] = t
			return m
		},
		ContentType: pkg.ContentTypeMultiPart,
	}
}

func init() {
	solutionCmd.AddCommand(listSolutionCmd)
	listSolutionCmd.Flags().BoolVarP(&createdOnly, "created", "c", false, "If true, limits the solutions to my (created) solutions.")
}
