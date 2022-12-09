package cmd

import (
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.ibm.com/skol/itzcli/pkg"
)

// workspaceCmd represents the workspace command
var workspaceCmd = &cobra.Command{
	Use:    "workspace [workspace name]",
	Short:  "Executes and interacts with different workspaces.",
	Long:   `Executes and interacts with different workspaces.`,
	Args:   cobra.ExactArgs(1),
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

func init() {
	RootCmd.AddCommand(workspaceCmd)
}
