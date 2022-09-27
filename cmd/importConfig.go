package cmd

import (
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.ibm.com/skol/atkcli/pkg"

	"github.com/spf13/cobra"
)

var ocpnowCfg string

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:    "import",
	Short:  "Imports configuration from various sources",
	Long:   `Imports configuration from various sources, including ocpnow`,
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debugf("Importing configuration from <%s>...", ocpnowCfg)
		return importOcpnowConfig(ocpnowCfg)
	},
}

func importOcpnowConfig(cfg string) error {
	cliCfg := viper.ConfigFileUsed()
	logger.Tracef("Adding contents of <%s> to configuration file <%s>...", cfg, cliCfg)
	err := pkg.AppendToFile(cfg, cliCfg)
	if err != nil {
		return err
	}
	logger.Infof("Succesfully imported ocpnow project configuration into %s.", cliCfg)
	return nil
}

func init() {
	configureCmd.AddCommand(importCmd)
	importCmd.Flags().StringVarP(&ocpnowCfg, "from-ocpnow-project", "p", "", "Specifies the project.yaml file created by ocpnow.")
}
