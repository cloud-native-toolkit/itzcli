package cmd

import (
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:    DeployAction,
	Short:  "Deploys a build in a cluster",
	Long:   "Deploys a build into an existing cluster",
	PreRun: SetLoggingLevel,
}

var deployBuildCmd = &cobra.Command{
	Use:    BuildResource,
	Short:  "Deploys the given build to the specified cluster",
	Long:   "Deploys the given build to the specified cluster",
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debug("Deploying your build...")
		return nil
	},
}

func init() {
	deployCmd.AddCommand(deployBuildCmd)

	rootCmd.AddCommand(deployCmd)
}
