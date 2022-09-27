package cmd

import (
	logger "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// importCmd represents the import command
var listConfigCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists configuration from various sources",
	Long:  `Lists configuration from various sources, including ocpnow`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debugf("Importing configuration from <%s>...", &ocpnowCfg)
	},
}

func init() {
	configureCmd.AddCommand(listConfigCmd)
}
