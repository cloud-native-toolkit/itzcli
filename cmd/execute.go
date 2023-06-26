package cmd

import (
	"fmt"

	"github.com/cloud-native-toolkit/itzcli/pkg"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var executeCmd = &cobra.Command{
	Use:    ExecuteAction,
	Short:  "Executes workspaces or pipelines in a cluster",
	Long:   "Executes workspaces or pipelines in a cluster",
	PreRun: SetLoggingLevel,
}

var executeWorkspaceCmd = &cobra.Command{
	Use:    WorkspaceResource,
	Short:  "Executes the given workspace",
	Long:   "Executes the given workspace",
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Infof("Running command: %s", args[0])
		workspace := args[0]
		key := pkg.Keyify(workspace)
		if viper.Get(pkg.FlattenCommandName(cmd, key)) == nil {
			return fmt.Errorf("workspace does not exist: %s", workspace)
		}
		return pkg.DoContainerizedStep(cmd, key, nil, nil)
	},
}

var executePipelineCmd = &cobra.Command{
	Use:    PipelineResource,
	Short:  "Executes the given pipeline",
	Long:   "Executes the given pipeline",
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debug("Executing your pipeline...")
		return nil
	},
}

func init() {
	executeCmd.AddCommand(executeWorkspaceCmd)
	executeCmd.AddCommand(executePipelineCmd)

	rootCmd.AddCommand(executeCmd)
}
