package cmd

import (
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.ibm.com/skol/atkcli/pkg"

	"github.com/spf13/cobra"
)

// importCmd represents the import command
var listConfigCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists configuration from various sources",
	Long:  `Lists configuration from various sources, including ocpnow`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ocpCfg := viper.GetStringSlice("ocpnow.configFiles")
		if len(ocpCfg) == 0 {
			logger.Warnf("No ocpnow configuration files found")
			return nil
		}
		logger.Debugf("Listing configuration from <%s>...", ocpCfg[0])
		project, err := pkg.LoadProject(ocpCfg[0])
		if err != nil {
			return err
		}
		err = project.Write(configureCmd.OutOrStdout())
		return err
	},
}

func init() {
	configureCmd.AddCommand(listConfigCmd)
}
