package cmd

import (
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

var cfgFile string
var verbose bool
var debug bool

var ITZVersionString = "No Version Provided"

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "itz",
	Short: fmt.Sprintf("IBM Technology Zone (ITZ) Command Line Interface (CLI), version %s", ITZVersionString),
	Long:  `IBM Technology Zone (ITZ) Command Line Interface (CLI)`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute(version string) {
	ITZVersionString = version
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.itz/cli-config.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Prints verbose messages")
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "X", false, "Prints trace messaging for debugging")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(filepath.Join(home, ".itz"))
		viper.SetConfigType("yaml")
		viper.SetConfigName("cli-config")
	}

	viper.SetEnvPrefix("ITZ")
	// Configure the key replacer so that environment variables in the form of
	// "ITZ_RESERVATIONS_API_TOKEN" will map to "reservations.api.token", because
	// remember that the ITZ_ prefix is configured by the SetEnvPrefix() function
	// above.
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv() // read in environment variables that match

	// Configure some logger stuff
	logger.SetOutput(RootCmd.ErrOrStderr())

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logger.Debugf("Using config file: %s", viper.ConfigFileUsed())
	}
}

func SetLoggingLevel(cmd *cobra.Command, args []string) {
	if debug {
		logger.SetLevel(logger.TraceLevel)
		logger.Trace("Trace logging enabled.")
		return
	}
	if verbose {
		logger.SetLevel(logger.DebugLevel)
		logger.Debug("Debug logging enabled.")
		return
	}
	// else, set it to warn only and format it a bit differently...
	logger.SetLevel(logger.InfoLevel)
}

func SetQuietLogging(cmd *cobra.Command, args []string) {
	logger.SetLevel(logger.WarnLevel)
}
