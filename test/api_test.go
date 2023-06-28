package test

// import (
// 	"fmt"
// 	"io"
// 	"io/ioutil"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/cloud-native-toolkit/itzcli/api"
// 	"github.com/cloud-native-toolkit/itzcli/cmd"
// 	"github.com/spf13/cobra"
// 	"github.com/stretchr/testify/assert"
// )

// func TestLogin(t *testing.T) {
// 	r := api.SetUpRouter(cmd.RootCmd)
// 	missingTokenreq, err := http.NewRequest("GET", "/login", nil)
// 	assert.NoError(t, err)
// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, missingTokenreq)
// 	assert.Equal(t, http.StatusUnauthorized, w.Code)

//     badAccessTokenreq, err := http.NewRequest("GET", "/login?token=12345", nil)
// 	assert.NoError(t, err)
// 	w = httptest.NewRecorder()
// 	r.ServeHTTP(w, badAccessTokenreq)
// 	assert.Equal(t, http.StatusBadRequest, w.Code)
// }
// // create a test command
// var testCMD = &cobra.Command{
// 	Use:   "unit",
// 	Short: "u",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		io.WriteString(cmd.OutOrStdout(), "Unit")
// 	},
// }
// // create a subTest command
// var subTestCommand = &cobra.Command{
// 	Use:   "test",
// 	Short: "t",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		io.WriteString(cmd.OutOrStdout(), " test are cool!")
// 	},
// }

// func TestAPIRegistered(t *testing.T) {
// 	cmd.RootCmd.AddCommand(testCMD)
// 	testCMD.AddCommand(subTestCommand)
// 	unit, unitError := httpRequest("/unit")
// 	assert.NoError(t, unitError)
// 	assert.Equal(t, "Unit", unit)
// 	test, testError := httpRequest("/unit/test")
// 	assert.NoError(t, testError)
// 	assert.Equal(t, " test are cool!", test)
// }

// func httpRequest(endpoint string) (string, error) {
// 	// run the command in the background
// 	r := api.SetUpRouter(cmd.RootCmd)

// 	// Create a new HTTP request.
// 	req, err := http.NewRequest("GET", endpoint, nil)
// 	if err != nil {
// 		return  "", fmt.Errorf("unexpected error: %v", err)
// 	}

// 	// Create a test HTTP server and send the request to it.
// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)

// 	// Verify that the HTTP response code is 200 OK.
// 	if w.Code != http.StatusOK {
// 		return "", fmt.Errorf("unexpected status code: %v", w.Code)
// 	}
// 	responseData, _ := ioutil.ReadAll(w.Body)
// 	return string(responseData), nil
// }
