package api

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cloud-native-toolkit/itzcli/pkg"
	"github.com/cloud-native-toolkit/itzcli/pkg/auth"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func StartServer(rootCmd *cobra.Command) *gin.Engine {
	r := SetUpRouter(rootCmd)
	srv := &http.Server{
		Addr:    "localhost:8795",
		Handler: r,
	}
	logger.Infof("Starting server on %s...", srv.Addr)
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			logger.Infof("listen: %s\n", err)
		}
	}()
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown:", err)
	}
	logger.Info("Server exiting")
	return r
}

func SetUpRouter(rootCmd *cobra.Command) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	gin.DefaultWriter = rootCmd.ErrOrStderr()
	r.GET("/login", GetTechZoneToken)
	CliToRESTHandler(r, rootCmd)
	return r
}

func GetTechZoneToken(c *gin.Context) {
	// Grab the access token
	accessToken := c.Query("token")
	if accessToken == "" {
		c.Data(http.StatusUnauthorized, "text/html; charset=utf-8", []byte(errorHTML))
		logger.Debug("Missing required access token...")
		auth.ErrorGettingToken()
		return
	}
	// Okay, we have the access token so let's make our API call to TechZone now
	accessTokenURL := fmt.Sprintf("https://auth.techzone.ibm.com/user?access_token=%s", accessToken)

	techZoneData, err := pkg.ReadHttpGetTWithFunc(accessTokenURL, "", func(code int) error {
		logger.Debugf("Handling HTTP return code %d...", code)
		return nil
	})
	if err != nil {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(errorHTML))
		auth.ErrorGettingToken()
		return
	}
	// parse the body and grab requestJson
	jsoner := auth.NewJsonReader()
	techZoneDataR := bytes.NewReader(techZoneData)
	requestJson, err := jsoner.Read(techZoneDataR)
	var buf bytes.Buffer
	if err != nil || requestJson.Token == "" {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(errorHTML))
		auth.ErrorGettingToken()
		return
	}
	htmlParser, _ := template.New("itzcliapi").Parse(succesHTML)
	htmlParser.Execute(&buf, requestJson)
	htmlString := buf.String()
	// Serve the HTML page
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlString))
	// Write the token to the config file
	_ = auth.SaveTokenToConfig(requestJson.Token)
}

// helper function to see if string exists in slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func CliToRESTHandler(router *gin.Engine, rootCmd *cobra.Command) {
	for _, cobraCmd := range rootCmd.Commands() {
		// What we want is `itz solution list --list-all` becomes
		// <url>/api/itz/solution/list&list-all=true
		baseCmdName := cobraCmd.Name()
		ignoredPaths := []string{"login", "completion", "execute", "generate"}
		if contains(ignoredPaths, baseCmdName) {
			continue
		}
		registerAPI(baseCmdName, router, cobraCmd, rootCmd)
		for _, subCommands := range cobraCmd.Commands() {
			subCommandName := subCommands.Use
			apiPath := fmt.Sprintf("/%s/%s", baseCmdName, subCommandName)
			registerAPI(apiPath, router, subCommands, rootCmd)
		}
	}
}

func registerAPI(path string, router *gin.Engine, command *cobra.Command, rootCmd *cobra.Command) {
	router.GET(path, func(c *gin.Context) {
		// Execute the command with the parsed flags and capture the output
		var args []string
		parent := command.Parent().Use
		if parent != "itz" {
			args = append(args, parent)
		}
		args = append(args, command.Use)
		args = append(args, "--json")
		// Parse the command name and flags from the HTTP request
		command.LocalFlags().VisitAll(func(flag *pflag.Flag) {
			value := c.Request.URL.Query().Get(flag.Name)
			if value != "" {
				flag.Value.Set(value)
				args = append(args, fmt.Sprintf("--%s=%s", flag.Name, value))
			}
		})
		rootCmd.SetArgs(args) // set the command's args
		rootCmd.SetOut(c.Writer)
		err := rootCmd.Execute()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		// Return the command output as an HTTP response
		c.Header("Content-Type", "application/json; charset=utf-8")
		// unset the values that were set
		command.LocalFlags().VisitAll(func(flag *pflag.Flag) {
			_ = flag.Value.Set("")
		})
	})
}
