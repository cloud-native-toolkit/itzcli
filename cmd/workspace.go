package cmd

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.ibm.com/skol/atkmod"
	"github.ibm.com/skol/itzcli/pkg"
	"strings"
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
		// First, look up the configuration for the command
		key := pkg.Keyify(workspace)
		serviceConfig := &pkg.ServiceConfig{}
		err := viper.UnmarshalKey(fmt.Sprintf("%s.%s", "workspaces", key), &serviceConfig, options)
		if err != nil {
		}
		logger.Tracef("Found configuration for key %s: %v", key, serviceConfig)

		logger.Infof("Starting workspace: %s", workspace)
		wsImg := &pkg.Service{
			CfgPrefix:   key,
			DisplayName: workspace,
			ImgName:     serviceConfig.Image,
			IsLocal:     serviceConfig.Local,
			Start:       pkg.RunHandler,
			MapToUID:    -1,
			VolumeOpt:   serviceConfig.MountOpts,
			Flags:       []string{"--rm", "-it"},
		}
		if len(serviceConfig.LocalDir) > 0 {
			if len(serviceConfig.RemoteDir) == 0 {
				logger.Warnf("remotedir not specified; defaulting to mounting to /workspace")
				serviceConfig.RemoteDir = "/workspace"
			}
			wsImg.Volumes = map[string]string{
				serviceConfig.LocalDir: serviceConfig.RemoteDir,
			}
		}
		for _, v := range serviceConfig.Volumes {
			if wsImg.Volumes == nil {
				wsImg.Volumes = make(map[string]string)
			}
			vmap := strings.Split(v, ":")
			if len(vmap) != 2 {
				return fmt.Errorf("invalid volume specification: %s", v)
			}
			wsImg.Volumes[vmap[0]] = vmap[1]
		}
		services := []pkg.Service{*wsImg}
		ctx := &atkmod.RunContext{
			Out: reservationCmd.OutOrStdout(),
			Err: reservationCmd.ErrOrStderr(),
			In:  reservationCmd.InOrStdin(),
		}
		err = pkg.StartupServices(ctx, services, pkg.Sequential)
		if err != nil {
			logger.Error("error is: %v", err)
		}
		return err
	},
}

func options(config *mapstructure.DecoderConfig) {
	config.ErrorUnused = false
	config.ErrorUnset = false
	config.IgnoreUntaggedFields = true
}

func init() {
	RootCmd.AddCommand(workspaceCmd)
}
