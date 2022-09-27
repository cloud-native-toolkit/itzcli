package cmd

import (
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.ibm.com/skol/atkcli/pkg"
	"path/filepath"
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
	// First thing we're going to do is to copy the file into the .atk home directory..
	// HACK: Here we should peek inside the file and use the project name as
	// the config file name.
	configDir := filepath.Dir(viper.ConfigFileUsed())
	importedCfg := filepath.Join(configDir, "project.yaml")
	err := pkg.AppendToFile(cfg, importedCfg)
	if err != nil {
		return err
	}
	logger.Tracef("Storing project config <%s> in configuration directory <%s>...", cfg, configDir)
	currentFiles := viper.GetStringSlice("ocpnow.configFiles")
	currentFiles = append(currentFiles, importedCfg)
	viper.Set("ocpnow.configFiles", currentFiles)
	viper.WriteConfig()
	return nil
}

func init() {
	configureCmd.AddCommand(importCmd)
	importCmd.Flags().StringVarP(&ocpnowCfg, "from-ocpnow-project", "p", "", "Specifies the project.yaml file created by ocpnow.")
}
