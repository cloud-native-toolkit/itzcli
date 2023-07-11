package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloud-native-toolkit/itzcli/pkg/configuration"

	"github.com/cloud-native-toolkit/itzcli/pkg"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var verbose bool
var debug bool
var jsonFormat bool

var ITZVersionString = "No Version Provided"

const TextCommandOutputFormat string = "text"
const JsonCommandOutputFormat string = "json"

// RootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:    "itz",
	Short:  fmt.Sprintf("IBM Technology Zone (ITZ) Command Line Interface (CLI), version %s", ITZVersionString),
	Long:   `IBM Technology Zone (ITZ) Command Line Interface (CLI)`,
	PreRun: SetLoggingLevel,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute(version string) {
	ITZVersionString = version
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.itz/cli-config.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Prints verbose messages")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "X", false, "Prints trace messaging for debugging")
	// changes the output format
	rootCmd.PersistentFlags().BoolVar(&jsonFormat, "json", false, "Changes output to JSON")
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
	logger.SetOutput(rootCmd.ErrOrStderr())

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logger.Debugf("Using config file: %s", viper.ConfigFileUsed())
	}
}

func SetLoggingLevel(cmd *cobra.Command, args []string) {
	if debug {
		logger.SetLevel(logger.TraceLevel)
		logger.Trace("Trace logging enabled...")
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

func LoadApiClientConfig(path string) (*configuration.ApiConfig, error) {
	cfg := &configuration.ServiceConfig{}

	err := viper.UnmarshalKey(path, &cfg, pkg.ConfigOptions)
	if err != nil {
		return nil, err
	}
	logger.Tracef("Found configuration for key %s: %v", path, cfg)
	return &cfg.API, nil
}

func GetFormat(cmd *cobra.Command) string {
	if jsonFormat {
		return JsonCommandOutputFormat
	}
	return TextCommandOutputFormat
}

func pluralOf(s string) string {
	return fmt.Sprintf("%s%s", s, "s")
}

type AssertMatcher func(string) bool

func NotNull(val string) bool {
	return len(strings.TrimSpace(val)) > 0
}

func ValidURL(val string) bool {
	_, err := url.ParseRequestURI(val)
	return NotNull(val) && err == nil
}

func AssertFlag(val string, match AssertMatcher, msg string) error {
	if !match(val) {
		return fmt.Errorf(msg)
	}
	return nil
}
