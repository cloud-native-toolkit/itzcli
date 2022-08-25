/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.ibm.com/Nathan-Good/atkcli/internal/outputters"
	"github.ibm.com/Nathan-Good/atkdep"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists the modules available",
	Long: `Lists the modules that can be installed by the Activation
ToolKit
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debug("Getting modules from the module index...")
		idxReader := atkdep.NewAtkGitHubIndexReader("test")
		entries, err := idxReader.List()
		if err != nil {
			return err
		}

		outputters.WriteEntries(entries, modulesCmd.OutOrStdout())
		return nil
	},
}

func init() {
	modulesCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
