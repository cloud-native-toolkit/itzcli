package cmd

import (
	"github.com/spf13/cobra"
)

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "List, configure, and manage clusters created by ocpnow",
}

func init() {
	RootCmd.AddCommand(clusterCmd)
}
