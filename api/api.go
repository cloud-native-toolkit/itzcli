package api

import (
	"bytes"
	"context"
	"fmt"
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

var RootCmd *cobra.Command

func StartServer(rootCmd *cobra.Command) *gin.Engine {
	r := SetUpRouter(rootCmd)
	srv := &http.Server{
        Addr:    "localhost:8795",
        Handler: r,
    }
	go func() {
        // service connections
        if err := srv.ListenAndServe(); err != nil {
            logger.Printf("listen: %s\n", err)
        }
    }()
	// Wait for interrupt signal to gracefully shutdown the server with
    // a timeout of 5 seconds.
    quit := make(chan os.Signal)
    signal.Notify(quit, os.Interrupt)
    <-quit
    logger.Println("Shutdown Server ...")

    ctx, cancel := context.WithTimeout(context.Background(), 60 * time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        logger.Fatal("Server Shutdown:", err)
    }
    logger.Println("Server exiting")
	return r
}

func SetUpRouter(rootCmd *cobra.Command) *gin.Engine {
	RootCmd = rootCmd
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	currentDir, _ := os.Getwd()
	r.LoadHTMLGlob(fmt.Sprintf("%s/templates/*", currentDir))
	r.GET("/login", GetTechZoneToken)
	CliToRESTHandler(r)
	return r
}

func GetTechZoneToken(c *gin.Context) {
	// Grab the access token
	accessToken := c.Query("token")
	if accessToken == "" {
		c.HTML(http.StatusUnauthorized, "error.html", "")
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
		c.HTML(http.StatusBadRequest, "error.html", "")
		auth.ErrorGettingToken()
		return
	}
	// parse the body and grab requestJson
	jsoner := auth.NewJsonReader()
	techZoneDataR := bytes.NewReader(techZoneData)
	requestJson, err := jsoner.Read(techZoneDataR)
	if err != nil || requestJson.Token == "" {
		c.HTML(http.StatusBadRequest, "error.html", "")
		auth.ErrorGettingToken()
		return
	}
	// Serve the HTML page
	c.HTML(http.StatusOK, "index.html", requestJson)
	// Write the token to the config file
	_ = auth.SaveTokenToConfig(requestJson.Token)
}

func CliToRESTHandler(router *gin.Engine) {
	for _, cobraCmd := range RootCmd.Commands() {
		// What we want is `itz solution list --list-all` becomes
		// <url>/api/itz/solution/list&list-all=true
		baseCmdName := cobraCmd.Name()
		if baseCmdName == "workspace" || baseCmdName == "completion" || baseCmdName == "api" {
			continue
		}
		registerAPI(baseCmdName, router, cobraCmd)
		for _, subCommands := range cobraCmd.Commands() {
			subCommandName := subCommands.Use
			apiPath := fmt.Sprintf("/%s/%s", baseCmdName, subCommandName)
			registerAPI(apiPath, router, subCommands)
		}
	}
}

func registerAPI(path string, router *gin.Engine, command *cobra.Command) {
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
		RootCmd.SetArgs(args) // set the command's args
		var buf bytes.Buffer
		RootCmd.SetOut(&buf)
		err := RootCmd.Execute()
		output := buf.String()
		httpStatus := http.StatusOK
		if err != nil {
			httpStatus = http.StatusInternalServerError
		}
		// Return the command output as an HTTP response
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.String(httpStatus, output)
		// unset the values that were set
		command.LocalFlags().VisitAll(func(flag *pflag.Flag) {
			_ = flag.Value.Set("")
		})
	})
}
